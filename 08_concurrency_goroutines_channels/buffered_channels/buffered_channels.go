package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. Unbuffered Channel
	// Sends block until there's a receiver.
	fmt.Println("--- Unbuffered Channel ---")
	unbuffered := make(chan string)
	go func() {
		fmt.Println("Sender: Sending message...")
		unbuffered <- "ping" // Blocks here until receiver is ready
		fmt.Println("Sender: Message sent")
	}()

	time.Sleep(500 * time.Millisecond) // Simulate some work
	fmt.Println("Receiver: Waiting...")
	msg := <-unbuffered
	fmt.Println("Receiver: Received", msg)

	// 2. Buffered Channel
	// Sends block only when the buffer is full.
	fmt.Println("\n--- Buffered Channel ---")
	buffered := make(chan string, 2) // Buffer size of 2

	fmt.Println("Sender: Sending 1...")
	buffered <- "one" // Won't block
	fmt.Println("Sender: Sending 2...")
	buffered <- "two" // Won't block

	// buffered <- "three" // This WOULD block because buffer is full

	fmt.Println("Receiver: Reading 1...", <-buffered)
	fmt.Println("Receiver: Reading 2...", <-buffered)
}
