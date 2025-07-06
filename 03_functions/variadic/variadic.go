// variadic.go
// Demonstrates variadic functions in Go.
package main

import "fmt"

// Variadic function
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println("Sum:", sum(1, 2, 3))
	fmt.Println("Sum with slice:", sum([]int{4, 5, 6}...))
}
// Documentation:
// - Variadic functions accept zero or more arguments.
// - Use ... to pass a slice as variadic arguments.
// - Edge cases: zero arguments, mixing variadic and non-variadic parameters.
