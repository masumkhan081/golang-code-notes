# 08_concurrency

Goroutines, channels, sync primitives, and concurrency patterns.
**Ordered easy → hard.** Read in order if you're learning; jump if you're revising.

## Contents

### Channels & goroutines — basics
- `01_goroutines_channels/` — spawning goroutines, send/receive on channels
- `02_buffered_channels/` — buffered vs unbuffered semantics
- `03_select/` — multi-channel multiplexing, `default` case
- `04_channel_closing/` — close convention, context-based cancellation
- `05_nil_channels/` — nil channel semantics (block forever, disable select cases)

### Synchronization
- `06_mutex_protected_cache/` — Mutex-guarded shared state
- `07_sync_primitives/` — `sync.Once`, `RWMutex`, `Cond`, `sync.Map`, `sync/atomic`

### Anti-patterns & detection
- `08_goroutine_leaks/` — leak patterns, fixes, and the `-race` detector

### Concurrency patterns
- `09_worker_pool/` — fixed pool consuming a job channel
- `10_fan_out_in/` — fan-out to N workers, fan-in via merge
- `11_pipeline/` — stage-based pipelines
- `12_semaphore_bounded_concurrency/` — buffered-channel semaphore
- `13_rate_limiter/` — token bucket via `golang.org/x/time/rate`
- `14_errgroup/` — `golang.org/x/sync/errgroup` fail-fast groups + bounded concurrency

## Cross-references

- Context cancellation: `../18_contexts_cancellation_timeout/`
- Interview Q&A: `../25_interview_questions/01_concurrency.md`
