package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID  int
	Output string
	Err    error
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			time.Sleep(100 * time.Millisecond)

			select {
			case <-ctx.Done():
				return
			case results <- Result{
				JobID:  job.ID,
				Output: fmt.Sprintf("worker-%d processed job-%d", id, job.ID),
			}:
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	workerCount := 3

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg)
	}

	go func() {
		defer close(jobs)
		for i := 1; i <= 8; i++ {
			jobs <- Job{ID: i}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result.Output)
	}
}
