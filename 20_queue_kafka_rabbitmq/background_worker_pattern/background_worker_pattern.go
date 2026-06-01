// background_worker_pattern.go
// Demonstrates a background worker that processes jobs from a channel,
// with graceful shutdown on OS signal. Common pattern in Go services
// for async tasks, queue consumers, and scheduled work.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Job represents a unit of work.
type Job struct {
	ID      int
	Payload string
}

// Worker processes jobs from jobCh until ctx is cancelled.
func worker(ctx context.Context, id int, jobCh <-chan Job, logger *slog.Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("worker shutting down", "worker_id", id)
			return
		case job, ok := <-jobCh:
			if !ok {
				logger.Info("job channel closed", "worker_id", id)
				return
			}
			processJob(ctx, id, job, logger)
		}
	}
}

func processJob(ctx context.Context, workerID int, job Job, logger *slog.Logger) {
	logger.Info("processing job", "worker_id", workerID, "job_id", job.ID, "payload", job.Payload)

	select {
	case <-ctx.Done():
		logger.Warn("job interrupted by cancellation", "job_id", job.ID)
		return
	case <-time.After(200 * time.Millisecond): // simulate work
		logger.Info("job done", "worker_id", workerID, "job_id", job.ID)
	}
}

// Dispatcher sends jobs into the channel and closes it when done.
func dispatcher(ctx context.Context, jobCh chan<- Job, total int, logger *slog.Logger) {
	defer close(jobCh) // signals workers that no more jobs will come

	for i := 1; i <= total; i++ {
		select {
		case <-ctx.Done():
			logger.Warn("dispatcher cancelled, sent only partial jobs", "sent", i-1)
			return
		case jobCh <- Job{ID: i, Payload: fmt.Sprintf("task-%d", i)}:
		}
	}
	logger.Info("dispatcher finished sending all jobs", "total", total)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx, cancel := context.WithCancel(context.Background())

	// Graceful shutdown on SIGINT / SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		logger.Info("shutdown signal received")
		cancel()
	}()

	const numWorkers = 3
	const numJobs = 10
	jobCh := make(chan Job, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			worker(ctx, i, jobCh, logger)
		}()
	}

	// Dispatch jobs
	go dispatcher(ctx, jobCh, numJobs, logger)

	// Wait for all workers to finish
	wg.Wait()
	logger.Info("all workers exited cleanly")
}
