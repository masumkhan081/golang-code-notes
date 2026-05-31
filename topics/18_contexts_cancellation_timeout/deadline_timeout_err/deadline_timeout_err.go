// deadline_timeout_err.go
// WithDeadline vs WithTimeout, plus ctx.Err() and ctx.Deadline() semantics.
//
// WithTimeout(parent, d) === WithDeadline(parent, time.Now().Add(d))
// They are not different features — WithTimeout is sugar.
//
// Pick by what you actually know:
//   WithTimeout  — "this operation must finish within X duration"
//   WithDeadline — "this operation must finish by absolute time T"
//                  (propagated deadlines, request budgets, scheduled jobs)
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func doWork(ctx context.Context, label string) error {
	// Report budget remaining (if any).
	if dl, ok := ctx.Deadline(); ok {
		fmt.Printf("[%s] budget=%v\n", label, time.Until(dl).Round(time.Millisecond))
	}

	select {
	case <-time.After(200 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err() // never returns nil after Done fires
	}
}

// Deadline propagation: child inherits the parent's deadline and may
// only shorten it, never extend. This is how request budgets cascade.
func deadlinePropagation() {
	parent, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Child gets at most 100ms — but bounded by parent's 500ms.
	child, cancel2 := context.WithTimeout(parent, 100*time.Millisecond)
	defer cancel2()

	if err := doWork(child, "child"); err != nil {
		fmt.Println("child:", err) // context deadline exceeded
	}

	// Parent still has ~400ms left.
	if err := doWork(parent, "parent"); err != nil {
		fmt.Println("parent:", err)
	} else {
		fmt.Println("parent: done")
	}
}

// Distinguishing cancellation from timeout via ctx.Err().
func classifyError(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, context.DeadlineExceeded):
		return "timeout"
	case errors.Is(err, context.Canceled):
		return "cancelled"
	default:
		return "other: " + err.Error()
	}
}

func main() {
	// 1. Plain timeout.
	{
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		err := doWork(ctx, "timeout")
		fmt.Println("→", classifyError(err))
	}

	// 2. Absolute deadline — equivalent form.
	{
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
		defer cancel()
		err := doWork(ctx, "deadline")
		fmt.Println("→", classifyError(err))
	}

	// 3. Explicit cancel.
	{
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(20 * time.Millisecond)
			cancel()
		}()
		err := doWork(ctx, "cancel")
		fmt.Println("→", classifyError(err))
	}

	// 4. Deadline propagation through call chain.
	deadlinePropagation()
}

/*
CHEATSHEET
  context.Background()         // root in main, init, tests
  context.TODO()               // placeholder while plumbing
  WithCancel(parent)           // returns ctx + cancel func — ALWAYS defer cancel
  WithTimeout(parent, d)       // sugar for WithDeadline(now+d)
  WithDeadline(parent, t)      // child deadline = min(parent deadline, t)
  WithValue(parent, key, val)  // request-scoped data only

ctx.Done()       — channel closed when ctx is cancelled or deadline hit
ctx.Err()        — nil before Done fires; Canceled or DeadlineExceeded after
ctx.Deadline()   — (time, true) if a deadline is set; (zero, false) otherwise
ctx.Value(key)   — request-scoped lookup; O(depth)

ALWAYS-DEFER-CANCEL
  Every WithCancel/WithTimeout/WithDeadline returns a cancel func.
  Forgetting defer cancel() leaks a goroutine + timer until the parent
  context is cancelled. `go vet` flags this with -vet=lostcancel.

CHECKING DEADLINE BEFORE STARTING
  if dl, ok := ctx.Deadline(); ok && time.Until(dl) < minNeeded {
      return context.DeadlineExceeded // fail fast, don't even try
  }
*/
