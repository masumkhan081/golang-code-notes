// errgroup.go
// golang.org/x/sync/errgroup — the modern idiom for "run N goroutines,
// fail fast on the first error, return when all are done".
//
// Replaces hand-written WaitGroup + error channel + mutex.
//
// Install:  go get golang.org/x/sync/errgroup
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

// fetchAll fans out one HTTP HEAD per URL, returns the first error,
// and cancels all sibling requests on failure.
func fetchAll(ctx context.Context, urls []string) error {
	// WithContext gives us a derived ctx that is cancelled the moment
	// any goroutine returns a non-nil error.
	g, ctx := errgroup.WithContext(ctx)

	for _, url := range urls {
		url := url // capture per iteration (pre-Go 1.22 idiom; harmless after)
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
			if err != nil {
				return err
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err // siblings get cancelled via ctx
			}
			resp.Body.Close()
			if resp.StatusCode >= 400 {
				return fmt.Errorf("%s: %s", url, resp.Status)
			}
			return nil
		})
	}

	return g.Wait() // first non-nil error, or nil
}

// Bounded concurrency: g.SetLimit(n) caps in-flight goroutines.
// Use this instead of a hand-rolled semaphore for most cases.
func bounded(ctx context.Context, jobs []int) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(4) // at most 4 goroutines at a time

	for _, j := range jobs {
		j := j
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(10 * time.Millisecond):
				fmt.Println("processed", j)
				return nil
			}
		})
	}
	return g.Wait()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := bounded(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8}); err != nil {
		fmt.Println("bounded error:", err)
	}

	if err := fetchAll(ctx, []string{
		"https://example.com",
		"https://example.org",
	}); err != nil {
		fmt.Println("fetchAll error:", err)
	}
}

/*
WHY ERRGROUP BEATS HAND-ROLLED WaitGroup
  - First error wins; siblings get cancelled via shared context.
  - No leaked error channels, no Mutex around an errs slice.
  - SetLimit gives you a built-in semaphore for bounded concurrency.

GOTCHAS
  - g.Wait() returns only the FIRST error. Sibling errors are dropped.
    If you need all errors, collect them yourself (errors.Join in Go 1.20+).
  - The derived ctx is cancelled on first error AND when Wait returns.
    Don't reuse that ctx after Wait.
  - g.Go is fire-and-forget — you cannot Wait then add more goroutines.

ALTERNATIVES
  - For "return all results" use a channel + WaitGroup (errgroup throws away
    results from later siblings).
  - For pipelines, errgroup pairs nicely — each stage = one g.Go.
*/
