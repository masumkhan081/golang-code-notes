# 12_context

Cancellation, timeouts, request-scoped values, graceful shutdown.

## Contents

- `graceful_shutdown.go` — HTTP server shutdown via context
- `context_misuse_pitfalls/` — anti-patterns (storing in struct, passing nil)
- `with_value_and_keys/with_value_and_keys.go` — **`WithValue` done right (typed keys, accessors, HTTP middleware example)**
- `deadline_timeout_err/deadline_timeout_err.go` — **`WithDeadline` vs `WithTimeout`, deadline propagation, classifying `ctx.Err()`**

## Cross-references

- Concurrency basics + patterns: `../11_concurrency/`
- errgroup (context-aware groups): `../11_concurrency/14_errgroup/`
- Interview Q&A: `../25_interview_questions/03_context.md`
