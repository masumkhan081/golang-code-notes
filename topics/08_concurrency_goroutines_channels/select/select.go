package main

import (
	"fmt"
	"time"
)

func main() {
	// Select allows a goroutine to wait on multiple communication operations.
	// It's like a switch statement, but for channels.

	c1 := make(chan string)
	c2 := make(chan string)

	// Simulate async operations
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// We'll loop twice to receive both values
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}

	// Non-blocking select with default
	// If no channel is ready, the default case is executed.
	c3 := make(chan string)
	select {
	case msg := <-c3:
		fmt.Println("received", msg)
	default:
		fmt.Println("no message received (non-blocking)")
	}
}
