package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	limiter := rate.NewLimiter(2, 2)

	ctx := context.Background()

	for i := 1; i <= 6; i++ {
		start := time.Now()
		if err := limiter.Wait(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("request %d allowed after %v\n", i, time.Since(start))
	}
}
