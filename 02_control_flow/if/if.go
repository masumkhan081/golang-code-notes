// if.go
// Demonstrates Go's if statement, including edge cases.
package main

import "fmt"

func main() {
	// Basic if
	n := 10
	if n > 5 {
		fmt.Println("n is greater than 5")
	}

	// if-else
	if n%2 == 0 {
		fmt.Println("n is even")
	} else {
		fmt.Println("n is odd")
	}

	// if with short statement
	if x := n * 2; x > 15 {
		fmt.Println("x is greater than 15:", x)
	}

	// Edge case: if with variable shadowing
	if n := 3; n < 5 {
		fmt.Println("shadowed n:", n)
	}
}
// Documentation:
// - if statements can include an initialization statement.
// - Variables declared in if's short statement are scoped to the if/else block.
// - Edge cases: shadowing, omitted else, single-line if.
