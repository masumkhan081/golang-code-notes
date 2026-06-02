package main

import (
    "fmt"
    "sync"
)

type Counter struct {
    mu     sync.Mutex
    counts map[string]int
}

func NewCounter() *Counter {
    return &Counter{
        counts: make(map[string]int),
    }
}

func (c *Counter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.counts[key]++
}

func (c *Counter) Snapshot() map[string]int {
    c.mu.Lock()
    defer c.mu.Unlock()

    out := make(map[string]int, len(c.counts))
    for k, v := range c.counts {
        out[k] = v
    }
    return out
}

func main() {
    counter := NewCounter()

    var wg sync.WaitGroup
    keys := []string{"task", "task", "email", "task", "email", "db"}

    for _, key := range keys {
        key := key
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := 0; i < 1000; i++ {
                counter.Inc(key)
            }
        }()
    }

    wg.Wait()
    fmt.Println(counter.Snapshot())
}
