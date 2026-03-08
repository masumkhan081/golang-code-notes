// for.go
// Demonstrates Go's for loop, including edge cases.
package main

import "fmt"

func main() {
	// Standard for loop
	for i := 0; i < 3; i++ {
		fmt.Println("i:", i)
	}

	// While-like loop
	n := 1
	for n < 5 {
		n *= 2
		fmt.Println("n:", n)
	}

	// Infinite loop (use with caution)
	count := 0
	for {
		count++
		if count > 2 {
			break
		}
		fmt.Println("infinite loop iteration:", count)
	}

	// Edge case: continue and break
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Println("odd i:", i)
	}
}
// Documentation:
// - Go's only loop keyword is for (no while/do-while).
// - Infinite loops are possible with for {}.
// - continue/break control loop flow.
// - Edge cases: loop variable scope, infinite loops, continue/break usage.
