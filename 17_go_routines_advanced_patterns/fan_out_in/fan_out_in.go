package main

import (
	"fmt"
	"sync"
	"time"
)

// Fan-Out: Multiple functions can read from the same channel until that channel is closed.
// Fan-In: A function can read from multiple inputs and proceed until all are closed.

func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func squareWorker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			// Simulate work
			time.Sleep(100 * time.Millisecond)
			out <- n * n
		}
		close(out)
	}()
	return out
}

// Fan-In: Merge multiple channels into one
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.
	// output copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are done.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	in := producer(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-Out: Distribute work across 2 workers
	c1 := squareWorker(in)
	c2 := squareWorker(in)

	// Fan-In: Consume the merged output from both workers
	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}
