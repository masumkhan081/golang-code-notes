// sync_primitives.go
// The standard sync toolbox beyond Mutex + WaitGroup.
//
// When to reach for which:
//   sync.Mutex      — exclusive access, mixed read/write
//   sync.RWMutex    — many readers, rare writers (read >> write)
//   sync.Once       — one-time initialization, lazy singletons
//   sync.Cond       — wait for a condition (rare; channels usually cleaner)
//   sync.Map        — concurrent map for stable-key, mostly-read workloads
//   sync/atomic     — lock-free counters and flags on primitive ints/pointers
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------------------
// 1. sync.Once — exactly-once initialization, even under concurrent calls.
// ---------------------------------------------------------------------------

type Config struct{ value string }

var (
	cfg     *Config
	cfgOnce sync.Once
)

func GetConfig() *Config {
	cfgOnce.Do(func() {
		// expensive: file load, network call, etc.
		cfg = &Config{value: "loaded"}
	})
	return cfg
}

// ---------------------------------------------------------------------------
// 2. sync.RWMutex — many readers, rare writers.
// RLock is reentrant-unsafe and DOES NOT upgrade to Lock — don't try.
// ---------------------------------------------------------------------------

type RWCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *RWCache) Get(k string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[k]
	return v, ok
}

func (c *RWCache) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[k] = v
}

// ---------------------------------------------------------------------------
// 3. sync/atomic — lock-free counters, flags, and pointer swaps.
// Use atomic.Int64 etc (Go 1.19+) instead of raw AddInt64 on *int64.
// ---------------------------------------------------------------------------

func atomicCounter() {
	var n atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n.Add(1)
		}()
	}
	wg.Wait()
	fmt.Println("atomic counter:", n.Load())
}

// atomic.Value / atomic.Pointer[T] for whole-struct swaps.
type Settings struct{ Timeout time.Duration }

var current atomic.Pointer[Settings]

func ReloadSettings(s *Settings) { current.Store(s) }
func CurrentSettings() *Settings { return current.Load() }

// ---------------------------------------------------------------------------
// 4. sync.Map — purpose-built concurrent map.
// Only beats RWMutex+map for: (a) keys written once, read many times,
// or (b) disjoint key sets per goroutine. Otherwise a plain map+mutex wins.
// ---------------------------------------------------------------------------

func syncMapDemo() {
	var m sync.Map
	m.Store("a", 1)
	m.Store("b", 2)

	if v, ok := m.Load("a"); ok {
		fmt.Println("sync.Map a =", v)
	}

	// LoadOrStore: atomic "get or set"
	actual, loaded := m.LoadOrStore("a", 99)
	fmt.Printf("a=%v loaded=%v\n", actual, loaded) // a=1 loaded=true

	m.Range(func(k, v any) bool {
		fmt.Printf("%v=%v\n", k, v)
		return true // false to stop
	})
}

// ---------------------------------------------------------------------------
// 5. sync.Cond — wait until a predicate becomes true.
// 95% of the time a channel is clearer. Use Cond only when you need
// broadcast-to-many-waiters or fine-grained state under a Mutex.
// ---------------------------------------------------------------------------

type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []int
}

func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Push(v int) {
	q.mu.Lock()
	q.items = append(q.items, v)
	q.mu.Unlock()
	q.cond.Signal() // wake one waiter; Broadcast() wakes all
}

func (q *Queue) Pop() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 { // ALWAYS loop — spurious wakeups exist
		q.cond.Wait() // atomically: unlock, sleep, re-lock on wake
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v
}

func main() {
	_ = GetConfig()
	_ = GetConfig() // initializer runs only once

	c := &RWCache{data: map[string]string{}}
	c.Set("k", "v")
	v, _ := c.Get("k")
	fmt.Println("cache:", v)

	atomicCounter()
	syncMapDemo()

	q := NewQueue()
	go func() {
		time.Sleep(50 * time.Millisecond)
		q.Push(42)
	}()
	fmt.Println("popped:", q.Pop())
}

/*
INTERVIEW TRAPS
  - sync.Mutex is NOT reentrant. Locking twice from the same goroutine
    deadlocks. Use a different design (split critical sections).
  - sync.WaitGroup.Add must be called BEFORE the goroutine starts,
    not inside it.
  - Don't copy sync types (Mutex, WaitGroup, Cond). `go vet` catches this.
    Pass *T or embed in a struct accessed by pointer.
  - sync.Map's Range iterates a snapshot — not consistent for concurrent writes.
  - atomic ops on non-atomic types via unsafe are undefined behavior.
    Stick to atomic.Int64 / atomic.Pointer[T] / sync/atomic functions.
*/
