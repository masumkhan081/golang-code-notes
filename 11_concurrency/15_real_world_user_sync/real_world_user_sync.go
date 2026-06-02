// real_world_user_sync.go
//
// Capstone example: a production-shaped concurrent API client.
//
// SCENARIO
//   We need to sync N user records from a remote API. Real-world constraints
//   we need to respect — none of which appear in a textbook "spawn 5 goroutines"
//   demo:
//
//     - rate limit       : the API allows ≤ 5 requests per second
//     - concurrency cap  : we don't want unbounded goroutines (memory + DOS risk)
//     - per-request timeout    : a hung server shouldn't wedge a worker
//     - overall deadline       : the whole job must finish within budget
//     - retry on transient err : 5xx and network errors get exponential backoff
//     - graceful shutdown      : Ctrl+C cancels in-flight work cleanly
//     - streaming results      : consumer processes results as they arrive
//     - observability          : per-request latency + final summary stats
//
// This file combines patterns from earlier sub-folders:
//   - worker pool          (09_worker_pool/)
//   - bounded concurrency  (12_semaphore_bounded_concurrency/)
//   - rate limiter         (13_rate_limiter/)
//   - channel pipelines    (11_pipeline/)
//   - context cancellation (../12_context/)
//   - leak-safe patterns   (08_goroutine_leaks/)
//
// Compare with the SEQUENTIAL run at the bottom to feel the win.
//
// Run:  go run real_world_user_sync.go
// API:  https://jsonplaceholder.typicode.com (free, no auth)
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/time/rate"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Job — one unit of work flowing through the pipeline.
type Job struct {
	UserID int
}

// Result — one outcome flowing out of the pipeline.
// We carry the request ID so we can correlate even on errors.
type Result struct {
	UserID   int
	User     *User
	Err      error
	Attempts int
	Latency  time.Duration
}

// ---------------------------------------------------------------------------
// Config — knobs you'd normally read from env/flags.
// ---------------------------------------------------------------------------

type Config struct {
	NumWorkers     int           // size of the worker pool
	RatePerSec     int           // tokens-per-second granted by the limiter
	RequestTimeout time.Duration // per single HTTP attempt
	OverallTimeout time.Duration // hard deadline for the whole job
	MaxRetries     int           // attempts beyond the first
	BackoffBase    time.Duration // first backoff sleep (doubles each retry)
}

var cfg = Config{
	NumWorkers:     3,
	RatePerSec:     5,
	RequestTimeout: 2 * time.Second,
	OverallTimeout: 30 * time.Second,
	MaxRetries:     3,
	BackoffBase:    200 * time.Millisecond,
}

// ---------------------------------------------------------------------------
// HTTP client — configured once, shared by all workers.
// http.DefaultClient has NO timeout. Never use it for outbound calls.
// ---------------------------------------------------------------------------

var httpClient = &http.Client{
	Timeout: cfg.RequestTimeout, // safety net; the per-request context is primary
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20, // reuse connections across workers to the same host
		IdleConnTimeout:     90 * time.Second,
	},
}

// ---------------------------------------------------------------------------
// fetchUser — a single attempt. Returns (user, isRetryable, error).
// ---------------------------------------------------------------------------

func fetchUser(ctx context.Context, id int) (*User, bool, error) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err // bad URL — never retry
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		// Network errors are transient — worth a retry unless ctx is done.
		retryable := !errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded)
		return nil, retryable, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 500:
		return nil, true, fmt.Errorf("server error: %s", resp.Status)
	case resp.StatusCode == 429:
		return nil, true, fmt.Errorf("rate limited by upstream: %s", resp.Status)
	case resp.StatusCode >= 400:
		return nil, false, fmt.Errorf("client error: %s", resp.Status) // 4xx: don't retry
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}

	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, false, fmt.Errorf("decode: %w", err) // malformed response — don't retry
	}
	return &u, false, nil
}

// ---------------------------------------------------------------------------
// fetchWithRetry — wraps fetchUser with backoff + jitter + context awareness.
// ---------------------------------------------------------------------------

func fetchWithRetry(ctx context.Context, limiter *rate.Limiter, id int) Result {
	start := time.Now()
	res := Result{UserID: id}

	for attempt := 1; attempt <= cfg.MaxRetries+1; attempt++ {
		res.Attempts = attempt

		// Rate limit — Wait blocks until a token is available OR ctx is cancelled.
		if err := limiter.Wait(ctx); err != nil {
			res.Err = fmt.Errorf("rate-limiter: %w", err)
			res.Latency = time.Since(start)
			return res
		}

		// Per-attempt timeout, derived from the parent.
		attemptCtx, cancel := context.WithTimeout(ctx, cfg.RequestTimeout)
		user, retryable, err := fetchUser(attemptCtx, id)
		cancel() // release the timer immediately, don't wait for defer

		if err == nil {
			res.User = user
			res.Latency = time.Since(start)
			return res
		}
		res.Err = err

		if !retryable || attempt > cfg.MaxRetries {
			res.Latency = time.Since(start)
			return res
		}

		// Exponential backoff with jitter — avoids thundering herd on shared failures.
		sleep := cfg.BackoffBase * (1 << (attempt - 1))
		jitter := time.Duration(rand.Int63n(int64(sleep / 2)))
		select {
		case <-time.After(sleep + jitter):
		case <-ctx.Done():
			res.Err = ctx.Err()
			res.Latency = time.Since(start)
			return res
		}
	}
	res.Latency = time.Since(start)
	return res
}

// ---------------------------------------------------------------------------
// worker — pulls jobs, sends results, exits when jobs channel closes
// or ctx is cancelled. Never blocks on send: we always select on ctx.Done().
// ---------------------------------------------------------------------------

func worker(ctx context.Context, id int, limiter *rate.Limiter, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return // jobs channel closed — no more work
			}
			r := fetchWithRetry(ctx, limiter, job.UserID)
			select {
			case results <- r:
			case <-ctx.Done():
				return // don't block on send if we're shutting down
			}
		}
	}
}

// ---------------------------------------------------------------------------
// runConcurrent — orchestrates the whole pipeline.
// ---------------------------------------------------------------------------

func runConcurrent(parentCtx context.Context, userIDs []int) []Result {
	ctx, cancel := context.WithTimeout(parentCtx, cfg.OverallTimeout)
	defer cancel()

	limiter := rate.NewLimiter(rate.Limit(cfg.RatePerSec), cfg.RatePerSec)

	jobs := make(chan Job)            // unbuffered: workers pull as ready
	results := make(chan Result, cfg.NumWorkers) // small buffer smooths bursts

	var wg sync.WaitGroup
	for i := 1; i <= cfg.NumWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, limiter, jobs, results, &wg)
	}

	// Producer: feed jobs, close channel when done. Runs in its own goroutine
	// so the main goroutine can stream results as they arrive.
	go func() {
		defer close(jobs)
		for _, id := range userIDs {
			select {
			case <-ctx.Done():
				return
			case jobs <- Job{UserID: id}:
			}
		}
	}()

	// Closer: when all workers finish, close results so the range below exits.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Drain — stream as they land. Real apps would write to a DB, send to
	// another channel, push to Kafka, etc.
	out := make([]Result, 0, len(userIDs))
	var ok, fail atomic.Int64
	for r := range results {
		if r.Err == nil {
			ok.Add(1)
			fmt.Printf("  ✓ #%d %-25s @%-12s %s  [%d attempt, %v]\n",
				r.User.ID, r.User.Name, r.User.Username, r.User.Email,
				r.Attempts, r.Latency.Round(time.Millisecond))
		} else {
			fail.Add(1)
			fmt.Printf("  ✗ #%d ERROR after %d attempt(s): %v\n", r.UserID, r.Attempts, r.Err)
		}
		out = append(out, r)
	}
	return out
}

// ---------------------------------------------------------------------------
// runSequential — same work, single goroutine, for comparison.
// ---------------------------------------------------------------------------

func runSequential(ctx context.Context, userIDs []int) []Result {
	out := make([]Result, 0, len(userIDs))
	limiter := rate.NewLimiter(rate.Limit(cfg.RatePerSec), cfg.RatePerSec)
	for _, id := range userIDs {
		r := fetchWithRetry(ctx, limiter, id)
		if r.Err == nil {
			fmt.Printf("  ✓ #%d %-25s @%-12s %s\n",
				r.User.ID, r.User.Name, r.User.Username, r.User.Email)
		} else {
			fmt.Printf("  ✗ #%d ERROR: %v\n", r.UserID, r.Err)
		}
		out = append(out, r)
	}
	return out
}

// ---------------------------------------------------------------------------
// stats — quick summary: success/fail counts + latency percentiles.
// ---------------------------------------------------------------------------

func summarize(label string, results []Result, total time.Duration) {
	var ok, fail int
	lats := make([]time.Duration, 0, len(results))
	for _, r := range results {
		if r.Err == nil {
			ok++
		} else {
			fail++
		}
		lats = append(lats, r.Latency)
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

	pct := func(p int) time.Duration {
		if len(lats) == 0 {
			return 0
		}
		i := (p * (len(lats) - 1)) / 100
		return lats[i].Round(time.Millisecond)
	}

	fmt.Printf("\n--- %s ---\n", label)
	fmt.Printf("  total time   : %v\n", total.Round(time.Millisecond))
	fmt.Printf("  success/fail : %d/%d\n", ok, fail)
	fmt.Printf("  latency p50  : %v\n", pct(50))
	fmt.Printf("  latency p95  : %v\n", pct(95))
	fmt.Printf("  latency max  : %v\n", pct(100))
}

// ---------------------------------------------------------------------------
// main — wires it all together, handles Ctrl+C for graceful shutdown.
// ---------------------------------------------------------------------------

func main() {
	// Go 1.20+ auto-seeds math/rand globally — no rand.Seed needed.

	// signal.NotifyContext returns a ctx cancelled on SIGINT/SIGTERM.
	// Workers and the producer all watch this ctx → Ctrl+C halts the pipeline
	// cleanly instead of killing the process mid-request.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	userIDs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Println("========== CONCURRENT (worker pool) ==========")
	t := time.Now()
	cResults := runConcurrent(ctx, userIDs)
	cTotal := time.Since(t)

	fmt.Println("\n========== SEQUENTIAL (baseline) ==========")
	t = time.Now()
	sResults := runSequential(ctx, userIDs)
	sTotal := time.Since(t)

	summarize("CONCURRENT", cResults, cTotal)
	summarize("SEQUENTIAL", sResults, sTotal)

	speedup := float64(sTotal) / float64(cTotal)
	fmt.Printf("\n⚡ Speedup: %.2fx — concurrency wins because requests overlap I/O wait.\n", speedup)
	fmt.Printf("   With workers=%d and rate=%d/s, the limiter is the bottleneck, not CPU.\n",
		cfg.NumWorkers, cfg.RatePerSec)
}

/*
WHAT THIS DEMONSTRATES (and where it lives in the repo)

  Worker pool                       — 09_worker_pool/
  Bounded concurrency               — 12_semaphore_bounded_concurrency/
  Rate limiting                     — 13_rate_limiter/
  Context cancellation + timeouts   — ../12_context/
  Leak-safe send (select on Done)   — 08_goroutine_leaks/
  Channel closing convention        — 04_channel_closing/
  Graceful shutdown via signal      — ../12_context/graceful_shutdown.go

DESIGN NOTES

  * Why a worker pool instead of "spawn one goroutine per request"?
    With 10 IDs it doesn't matter. With 10,000 IDs you'd open 10,000
    sockets, exhaust file descriptors, and DOS the upstream. The pool
    bounds resource usage to a known constant.

  * Why a rate limiter on top of the worker cap?
    Workers control YOUR concurrency. The rate limiter respects the
    UPSTREAM's contract. Both are necessary — they constrain different
    resources.

  * Why per-request timeout AND overall timeout?
    Per-request stops one slow server from blocking a worker forever.
    Overall stops the whole job from running past its budget when many
    requests are slow but not slow enough to time out individually.

  * Why classify errors retryable vs not?
    Retrying a 404 is pointless and wastes the rate budget. Retrying a
    5xx or a network blip is exactly what you want. Smart retries are
    cheap, dumb retries amplify outages.

  * Why exponential backoff with jitter?
    Without backoff: a flapping upstream gets pounded the instant it
    recovers. Without jitter: all clients retry in lockstep, creating
    a thundering herd. Both matter at scale.

  * Why a closer goroutine (`wg.Wait(); close(results)`)?
    The main goroutine ranges over results. Without close, the range
    blocks forever once workers finish. close MUST come from the side
    that knows all sends are done — here, after WaitGroup drains.

  * What would change for 1M users?
    - read IDs from a stream (channel/scanner), don't preload a slice
    - persist progress (resume on crash)
    - shard rate limiter per upstream host
    - publish results to Kafka instead of stdout
    - emit metrics (Prometheus counters/histograms) instead of printf

POSSIBLE EXTENSIONS

  * Fan-out a second stage (e.g., fetch each user's /posts) — that's a
    pipeline. See 11_pipeline/.
  * Use errgroup.WithContext + g.SetLimit(N) for an even tighter version.
    See 14_errgroup/. The hand-rolled version here is more explicit and
    easier to instrument with stats.
*/
