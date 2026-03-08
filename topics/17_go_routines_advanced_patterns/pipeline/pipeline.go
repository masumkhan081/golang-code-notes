package main

import (
	"fmt"
)

// Pipeline Pattern
// A series of stages connected by channels, where each stage is a group of goroutines.

// Stage 1: Generator - converts a list of integers to a channel
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// Stage 2: Square - receives integers, squares them, and sends them out
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// Set up the pipeline
	// gen() -> sq() -> main()

	// 1. Generate numbers
	c := gen(2, 3, 4, 5)

	// 2. Square numbers
	out := sq(c)

	// 3. Consume results
	fmt.Println("Pipeline results:")
	for n := range out {
		fmt.Println(n)
	}
}
