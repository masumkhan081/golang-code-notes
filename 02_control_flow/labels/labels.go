// labels.go
// Demonstrates Go's label and goto/break/continue usage.
package main

import "fmt"

func main() {
	// Basic label and goto
	count := 0
Loop:
	for {
		count++
		if count < 3 {
			fmt.Println("Looping", count)
			continue Loop
		}
		break Loop
	}
	fmt.Println("Exited loop at count:", count)

	// Nested loops with break/continue label
	outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == j {
				continue outer
			}
			fmt.Printf("i=%d, j=%d\n", i, j)
		}
	}
}
// Documentation:
// - Labels allow breaking/continuing outer loops.
// - Goto is rarely used in modern Go code.
// - Edge cases: label shadowing, misuse can reduce code clarity.
