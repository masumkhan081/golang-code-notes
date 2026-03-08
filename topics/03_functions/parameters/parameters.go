// parameters.go
// Demonstrates function parameters in Go, including value and pointer semantics.
package main

import "fmt"

// Value parameter
func add(a int, b int) int {
	return a + b
}

// Pointer parameter
func increment(x *int) {
	*x = *x + 1
}

func main() {
	// Value parameter
	sum := add(2, 3)
	fmt.Println("Sum:", sum)

	// Pointer parameter
	n := 10
	increment(&n)
	fmt.Println("Incremented n:", n)
}
// Documentation:
// - Function parameters are passed by value by default.
// - Use pointers to modify arguments in-place.
// - Edge cases: nil pointer dereference, type mismatch.
