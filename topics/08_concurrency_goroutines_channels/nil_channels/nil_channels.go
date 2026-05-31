// nil_channels.go
// Nil channel semantics — a frequent interview probe.
//
// Operation         | nil channel
// ------------------+------------------
// send (ch <- x)    | blocks FOREVER
// recv (<-ch)       | blocks FOREVER
// close(ch)         | PANICS
//
// Why useful: a nil case in select is effectively "disabled".
// You can dynamically turn select branches on and off by setting
// channels to nil instead of writing complex branching logic.
package main

import "fmt"

// merge reads from in1 and in2 until both are drained.
// When a channel is closed, we set its local var to nil to disable
// that case in the select — otherwise we'd spin reading the zero value.
func merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for in1 != nil || in2 != nil {
			select {
			case v, ok := <-in1:
				if !ok {
					in1 = nil // disable this case
					continue
				}
				out <- v
			case v, ok := <-in2:
				if !ok {
					in2 = nil // disable this case
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

func producer(values ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()
	return ch
}

func main() {
	a := producer(1, 3, 5)
	b := producer(2, 4, 6)
	for v := range merge(a, b) {
		fmt.Println(v)
	}

	// Demonstrate the block: uncomment to deadlock.
	// var ch chan int     // nil
	// ch <- 1             // fatal error: all goroutines are asleep - deadlock!

	// Demonstrate the panic: uncomment.
	// var ch chan int
	// close(ch) // panic: close of nil channel
}

/*
CHANNEL CHEATSHEET (memorize for interviews)

                     | nil       | open empty | open with data | closed
  ch <- v            | block     | block      | maybe block    | PANIC
  <-ch               | block     | block      | receive        | zero, ok=false
  close(ch)          | PANIC     | ok         | ok             | PANIC

OTHER RULES
  - Close from sender side only. Never close from receiver.
  - Don't close a channel with multiple senders without coordination
    (use sync.Once or a separate done channel).
  - len(ch) and cap(ch) work but are racy — useful only for metrics.
  - Sending on a closed channel panics — but receiving returns
    (zero value, false) forever; safe to range over.
*/
