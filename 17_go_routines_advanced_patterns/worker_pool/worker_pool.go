package main

import (
	"fmt"
	"time"
)

// Worker Pool Pattern
// Use this to limit the number of goroutines processing tasks.

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("worker %d started job %d\n", id, j)
		time.Sleep(time.Second) // Simulate expensive task
		fmt.Printf("worker %d finished job %d\n", id, j)
		results <- j * 2
	}
}

func main() {
	const numJobs = 5
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// 1. Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// 2. Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Close jobs channel to signal no more work

	// 3. Collect results
	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
