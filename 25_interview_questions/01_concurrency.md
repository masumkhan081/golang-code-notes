# Concurrency — goroutines, channels, sync

Cross-references go to working code in this repo.

---

### Q1. What's the difference between a goroutine and an OS thread?

A goroutine is a user-space coroutine managed by Go's runtime scheduler.
Initial stack is **~2 KB** (grows/shrinks dynamically) vs **1-2 MB** for an
OS thread. Goroutines are multiplexed onto a small pool of OS threads
(default `GOMAXPROCS` = CPU cores) using the **GMP scheduler**:
- **G** — goroutine
- **M** — OS thread (machine)
- **P** — logical processor (holds run queue, executes Gs on Ms)

When a G blocks on syscall/IO, the runtime parks it and the M either picks
another G from its P's queue or hands the P to a different M.

---

### Q2. Unbuffered vs buffered channel — when do you use each?

| | Unbuffered | Buffered |
|---|---|---|
| Send blocks until... | a receiver is ready | buffer has room |
| Use for | **synchronization** (hand-off, signal) | **decoupling** rate of producer/consumer |
| Default | yes — prefer this | only when buffer size is justified by domain |

Picking buffer size: 0 (unbuffered), 1 (signal/latch), or N where N matches
a known batch/burst. Arbitrary sizes (e.g. 100) hide bugs by masking back-pressure.

---

### Q3. Who closes a channel?

**The sender**, and only the sender. Closing from the receiver risks a
"send on closed channel" panic. With multiple senders, coordinate via
`sync.Once` or a separate `done` channel — never close from each sender.

Receiver-side detection:
```go
v, ok := <-ch
if !ok { /* channel closed and drained */ }
```

---

### Q4. What happens on a nil channel?

| Op | Result |
|---|---|
| `ch <- v` on nil | blocks forever |
| `<-ch` on nil | blocks forever |
| `close(ch)` on nil | panics |

This is useful: setting a channel to `nil` inside a `select` **disables**
that case dynamically. See `08_concurrency_goroutines_channels/nil_channels/`.

---

### Q5. What is a goroutine leak? Give three causes.

A goroutine that never exits. Holds memory and captured variables until
the process dies.

Causes:
1. **Sending on an unbuffered channel with no receiver** — sender blocks forever.
2. **Receiving from a channel that's never closed** — `for range ch` never returns.
3. **`select` without a cancellation case** — no `<-ctx.Done()`, so if none of the other cases fire, the goroutine wedges.

Detection: `runtime.NumGoroutine()`, `net/http/pprof` goroutine dump,
`github.com/uber-go/goleak`. Code: `08_concurrency_goroutines_channels/goroutine_leaks/`.

---

### Q6. Mutex vs RWMutex vs atomic — when?

- `sync.Mutex` — exclusive access. Default choice when you're unsure.
- `sync.RWMutex` — many readers, **rare** writers. Reader/writer book-keeping has overhead, so it loses to plain Mutex when writes are frequent or critical sections are tiny.
- `sync/atomic` — single primitive value (counter, flag, pointer). Lock-free, ~10x faster than Mutex for the right shape. Use `atomic.Int64`, `atomic.Pointer[T]`.

Rule: profile before reaching for RWMutex. Most code that "should be"
RWMutex is fine on Mutex.

---

### Q7. What's wrong with this code?

```go
for i := 0; i < 10; i++ {
    go func() { fmt.Println(i) }()
}
```

Before Go 1.22, all goroutines captured the **same** `i` variable — the
loop usually finished before they ran, so they all printed 10.

In Go 1.22+, each iteration gets a fresh `i`, so it prints 0-9 in some order.

Pre-1.22 fix:
```go
for i := 0; i < 10; i++ {
    i := i // shadow
    go func() { fmt.Println(i) }()
}
```

Or pass it as an argument:
```go
go func(i int) { fmt.Println(i) }(i)
```

Know both behaviors — the question filters seniority.

---

### Q8. What does `select` with `default` do?

Makes the select **non-blocking**. If no other case is immediately ready,
`default` runs.

```go
select {
case msg := <-ch:
    handle(msg)
default:
    // nothing waiting, move on
}
```

Used for: non-blocking sends/receives, polling, "try-lock" patterns.
Don't use it in a tight loop — that's a busy-spin.

---

### Q9. Why does `go vet` warn about copying a `sync.Mutex`?

A `Mutex`'s zero value is a valid unlocked mutex. Copying a locked one
gives you two structures that both think they hold the lock, leading to
double-unlock panics or lost mutual exclusion. Same for `WaitGroup`, `Cond`,
`atomic.Int64`. **Pass them by pointer, embed them in structs accessed by pointer.**

---

### Q10. What's a data race vs a race condition?

- **Data race** — two goroutines access the same memory, one is a write, no synchronization. Undefined behavior. The race detector (`-race`) finds these.
- **Race condition** — logic bug where ordering between operations matters. May or may not involve a data race. The detector does NOT find these.

Example race condition without a data race:
```go
mu.Lock(); ok := exists(k); mu.Unlock()    // synchronized read
if !ok {
    mu.Lock(); insert(k); mu.Unlock()      // synchronized write
}
// Another goroutine could insert between Unlock and Lock — TOCTOU.
```

Fix: hold the lock across the check-and-act, or use atomic LoadOrStore.

---

### Q11. Why is `time.After` in a `select` loop a leak?

```go
for {
    select {
    case msg := <-ch:
        handle(msg)
    case <-time.After(5 * time.Second): // BUG
        return errTimeout
    }
}
```

Every loop iteration creates a **new** timer. The garbage collector cannot
collect a timer until it fires — so a busy loop that gets messages every
millisecond allocates and pins thousands of unfired 5-second timers.

Fix — create the timer once and reset it:

```go
t := time.NewTimer(5 * time.Second)
defer t.Stop()
for {
    select {
    case msg := <-ch:
        handle(msg)
        if !t.Stop() { <-t.C }
        t.Reset(5 * time.Second)
    case <-t.C:
        return errTimeout
    }
}
```

Or use `context.WithTimeout` outside the loop. Go 1.23+ improved
`time.After` to not leak as badly, but the idiom is still discouraged.

---

### Q12. What does a receive from a closed channel return?

The zero value, with `ok = false` — **forever**, never blocks:

```go
close(ch)
v, ok := <-ch  // v = zero, ok = false
v, ok = <-ch   // same; safe to call repeatedly
```

This is why `for v := range ch` exits cleanly when the channel is closed.

Send on a closed channel **panics** instead. Closing an already-closed
channel also panics. Closing a `nil` channel panics. These three panics
are the bulk of channel bugs in production.

---

### Q13. Explain `errgroup` and when you'd use it.

`golang.org/x/sync/errgroup` runs N goroutines, returns the **first** error,
and cancels siblings via a shared context. Replaces hand-rolled
WaitGroup + error channel.

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(8)              // bounded concurrency
for _, url := range urls {
    url := url
    g.Go(func() error { return fetch(ctx, url) })
}
return g.Wait()
```

Caveat: only the first error is returned — sibling errors are dropped.
Code: `17_go_routines_advanced_patterns/errgroup/`.
