# 17_go_routines_advanced_patterns

Concurrency patterns and the standard sync toolbox.

## Contents

- `worker_pool/` — fixed pool of workers consuming a job channel
- `fan_out_in/` — fan-out to N workers, fan-in via merge
- `pipeline/` — stage-based pipelines with channels
- `semaphore_bounded_concurrency/` — buffered-channel semaphore
- `rate_limiter/` — `golang.org/x/time/rate` token bucket
- `mutex_protected_cache/` — Mutex-guarded shared state
- `sync_primitives/sync_primitives.go` — **`sync.Once`, `RWMutex`, `Cond`, `sync.Map`, `sync/atomic`**
- `errgroup/errgroup.go` — **`golang.org/x/sync/errgroup` for fail-fast goroutine groups + bounded concurrency**

## Cross-references

- Basics + leaks + race detector: `../08_concurrency_goroutines_channels/`
- Context cancellation: `../18_contexts_cancellation_timeout/`
- Interview Q&A: `../25_interview_questions/01_concurrency.md`
