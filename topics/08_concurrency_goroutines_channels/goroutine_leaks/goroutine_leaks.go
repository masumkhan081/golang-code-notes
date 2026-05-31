// goroutine_leaks.go
// Common goroutine leak patterns and how to fix them.
// A leaked goroutine never exits — it holds memory, file descriptors,
// and any captured variables until the process dies.
package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// LEAK #1: send on an unbuffered channel with no receiver.
// The producer blocks forever once the consumer returns early.
func leakSendBlocks() {
	ch := make(chan int) // unbuffered
	go func() {
		ch <- 42 // blocks forever — nobody reads
	}()
	// consumer never reads — goroutine leaks
}

// FIX #1: use a buffered channel sized for the send, or use select+context.
func fixedSendBlocks(ctx context.Context) {
	ch := make(chan int, 1) // buffer absorbs the send
	go func() {
		select {
		case ch <- 42:
		case <-ctx.Done():
			return
		}
	}()
}

// LEAK #2: receive from a channel nobody will close.
func leakReceiveBlocks() {
	ch := make(chan int)
	go func() {
		for v := range ch { // blocks forever if sender never closes
			fmt.Println(v)
		}
	}()
}

// FIX #2: sender owns the close. Always.
func fixedReceiveBlocks() {
	ch := make(chan int)
	go func() {
		defer close(ch) // sender closes
		for i := 0; i < 3; i++ {
			ch <- i
		}
	}()
	for v := range ch {
		fmt.Println(v)
	}
}

// LEAK #3: select with no default and no cancellation.
// If none of the cases ever fire, the goroutine wedges.
func leakSelectNoExit(ch1, ch2 chan int) {
	go func() {
		select {
		case <-ch1:
		case <-ch2:
		}
	}()
}

// FIX #3: always include a cancellation case (context).
func fixedSelectExit(ctx context.Context, ch1, ch2 chan int) {
	go func() {
		select {
		case <-ch1:
		case <-ch2:
		case <-ctx.Done():
			return
		}
	}()
}

// LEAK #4: WaitGroup.Add called after Wait, or Done forgotten.
// Wait blocks forever; the calling goroutine leaks too.
func leakWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// forgot wg.Done() — Wait blocks forever
	}()
	wg.Wait()
}

// FIX #4: defer wg.Done() at the top of the goroutine.
func fixedWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// work
	}()
	wg.Wait()
}

// LEAK #5: HTTP / DB call without context timeout.
// If the server hangs, the goroutine making the call hangs forever.
// FIX: always pass a context with a timeout (see 18_contexts).

// Detecting leaks at runtime: runtime.NumGoroutine() before/after a test.
// In real code use github.com/uber-go/goleak in TestMain.
func goroutineCount() int {
	return runtime.NumGoroutine()
}

// ----------------------------------------------------------------------------
// RACE CONDITIONS — run examples with: go run -race goroutine_leaks.go
// ----------------------------------------------------------------------------

// RACE: concurrent write without sync. Detector flags this.
func racyCounter() {
	var counter int
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // DATA RACE
		}()
	}
	wg.Wait()
	fmt.Println("racy:", counter) // value is non-deterministic
}

// FIX with mutex.
func safeCounterMutex() {
	var (
		mu      sync.Mutex
		counter int
		wg      sync.WaitGroup
	)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("safe:", counter)
}

func main() {
	before := goroutineCount()

	fixedReceiveBlocks()
	fixedWaitGroup()
	safeCounterMutex()

	// Give finished goroutines time to be reaped before measuring.
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("goroutines: before=%d after=%d\n", before, goroutineCount())
}

/*
KEY RULES
  1. The sender closes the channel, never the receiver.
  2. Every goroutine needs a guaranteed exit path — usually <-ctx.Done().
  3. Buffer a channel only when you know the producer/consumer ratio.
  4. defer wg.Done() at the top of the goroutine, not the bottom.
  5. Run tests with -race in CI. Always.

RACE DETECTOR
  go test -race ./...
  go run -race main.go
  go build -race
  Cost: ~2-20x slower, ~5-10x more memory. Use in CI, not prod binaries.

DETECTION TOOLS
  - runtime.NumGoroutine()        // cheap sanity check
  - net/http/pprof goroutine dump // production diagnosis
  - github.com/uber-go/goleak     // leak assertions in tests
*/
