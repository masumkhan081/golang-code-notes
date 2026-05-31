# 08_concurrency_goroutines_channels

Core Go superpower — parallelism.

## Contents

- `goroutines_channels/` — spawn goroutines, send/receive on channels
- `buffered_channels/` — buffered vs unbuffered semantics
- `select/` — multi-channel multiplexing, default case
- `channel_closing_and_cancellation.go` — close convention, context-based cancellation
- `nil_channels/nil_channels.go` — **nil channel semantics; disabling select cases dynamically**
- `goroutine_leaks/goroutine_leaks.go` — **leak patterns, fixes, and the `-race` detector**

## Cross-references

- Advanced patterns (worker pool, fan-in/out, pipeline, errgroup, sync primitives): `../17_go_routines_advanced_patterns/`
- Context cancellation: `../18_contexts_cancellation_timeout/`
- Interview Q&A: `../25_interview_questions/01_concurrency.md`
