package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const maxConcurrent = 3
	sem := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		i := i
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Println("starting job", i)
			time.Sleep(300 * time.Millisecond)
			fmt.Println("finished job", i)
		}()
	}

	wg.Wait()
}
