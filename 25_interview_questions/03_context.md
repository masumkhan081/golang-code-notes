# Context

---

### Q1. What is `context.Context` for?

Three jobs, in priority order:
1. **Cancellation propagation** — tell goroutines to stop work.
2. **Deadlines/timeouts** — bound how long an operation can take.
3. **Request-scoped values** — IDs, auth, tracing data crossing API boundaries.

It is **not** for dependency injection, configuration, or anything that
isn't request-scoped.

---

### Q2. How do you pass a context?

- Always **first argument**, named `ctx`.
- Never store it in a struct field. Pass it through call chains.
- Never pass `nil` — use `context.Background()` at roots (main, init, tests) or `context.TODO()` while plumbing is incomplete.

---

### Q3. `WithTimeout` vs `WithDeadline`?

`WithTimeout(parent, d)` is sugar for `WithDeadline(parent, time.Now().Add(d))`.
Same primitive. Pick by what you're expressing:
- duration budget → `WithTimeout`
- absolute time / propagated upstream deadline → `WithDeadline`

A child's deadline is `min(parent, child)` — children can shorten, never extend.

---

### Q4. What does `ctx.Done()` return and when does it fire?

A `<-chan struct{}` that is **closed** when the context is cancelled
(explicit cancel, timeout, or parent cancellation). Closed channels never
re-open, so once Done fires, any future receive returns immediately.

```go
select {
case <-time.After(d):
    // work done
case <-ctx.Done():
    return ctx.Err() // never nil after Done fires
}
```

---

### Q5. How do you distinguish a timeout from an explicit cancel?

```go
switch {
case errors.Is(err, context.DeadlineExceeded):
    // timeout
case errors.Is(err, context.Canceled):
    // explicit cancel() call
}
```

`ctx.Err()` returns one of these after `Done()` fires; nil before.

---

### Q6. Why is `defer cancel()` mandatory?

Every `WithCancel`/`WithTimeout`/`WithDeadline` returns a cancel function.
Even if the operation finishes early, calling cancel releases the timer
and unblocks resources tied to the context. Forgetting it leaks a small
amount of memory + a goroutine per leak until the parent context dies.

`go vet -vet=lostcancel` warns about missing defers.

---

### Q7. How should you use `WithValue`?

For **request-scoped** data only (request ID, auth user, trace span).

Rules:
- Key must be an **unexported custom type** to prevent cross-package collisions.
- Expose typed accessor functions; never the key.
- Don't carry mutable state or pointers to mutable state — that turns ctx into a back-door global.

```go
type ctxKey int
const keyRequestID ctxKey = iota

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, keyRequestID, id)
}
func RequestID(ctx context.Context) (string, bool) {
    v, ok := ctx.Value(keyRequestID).(string)
    return v, ok
}
```

Code: `12_context/with_value_and_keys/`.

---

### Q8. What's wrong with storing ctx in a struct?

- Lifetime gets confused — the struct outlives the request that created the ctx.
- Goroutines created later may capture a long-cancelled ctx, or a request-scoped ctx after the request returned.
- It hides the parameter — readers can't see what's cancellable.

Pass ctx as a function argument every time.

---

### Q9. What's the performance cost of `ctx.Value`?

`Value` walks a linked list of parent contexts — **O(depth)**. Fine for
typical 3-5 deep chains. Don't put hot-path data there; don't make it your
primary lookup mechanism.

---

### Q10. How do you check budget before doing expensive work?

```go
if dl, ok := ctx.Deadline(); ok && time.Until(dl) < estimatedCost {
    return context.DeadlineExceeded // fail fast
}
```

Useful before starting a long DB query or network call you know won't fit
in the remaining budget.
