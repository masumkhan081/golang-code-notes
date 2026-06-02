package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

func producer(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case <-ctx.Done():
                return
            case out <- n:
            }
        }
    }()
    return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                return
            case n, ok := <-in:
                if !ok {
                    return
                }
                time.Sleep(100 * time.Millisecond)
                select {
                case <-ctx.Done():
                    return
                case out <- n * n:
                }
            }
        }
    }()
    return out
}

func merge(ctx context.Context, chans ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(chans))

    for _, ch := range chans {
        ch := ch
        go func() {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case n, ok := <-ch:
                    if !ok {
                        return
                    }
                    select {
                    case <-ctx.Done():
                        return
                    case out <- n:
                    }
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    in := producer(ctx, 1, 2, 3, 4, 5, 6)
    c1 := square(ctx, in)
    c2 := square(ctx, in)

    count := 0
    for n := range merge(ctx, c1, c2) {
        fmt.Println("result:", n)
        count++
        if count == 3 {
            cancel()
        }
    }

    fmt.Println("done without leaking blocked goroutines")
}
