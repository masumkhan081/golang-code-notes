// goroutines_channels.go
// Demonstrates goroutines, channels, and select in Go.
package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("worker %d started job %d\n", id, j)
		time.Sleep(time.Second)
		fmt.Printf("worker %d finished job %d\n", id, j)
		results <- j * 2
	}
}

func main() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for w := 1; w <= 2; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 5; a++ {
		fmt.Println("result:", <-results)
	}

	// Edge case: select with timeout
	c := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		c <- "done"
	}()
	select {
	case msg := <-c:
		fmt.Println(msg)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}
// Documentation:
// - Goroutines run functions concurrently.
// - Channels are used for communication between goroutines.
// - Select allows waiting on multiple channel operations.
// - Edge cases: deadlocks, race conditions, channel closing.
