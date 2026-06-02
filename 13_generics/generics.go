// generics.go
// Demonstrates Go generics (Go 1.18+).
package main

import "fmt"

// Generic function
func PrintSlice[T any](s []T) {
	for _, v := range s {
		fmt.Println(v)
	}
}

// Generic type
type Pair[T, U any] struct {
	First  T
	Second U
}

func main() {
	PrintSlice([]int{1, 2, 3})
	PrintSlice([]string{"a", "b", "c"})

	p := Pair[int, string]{First: 1, Second: "one"}
	fmt.Println("Pair:", p)
}
// Documentation:
// - Generics use type parameters in [] after function/type name.
// - any is a type constraint for any type.
// - Edge cases: type inference, constraints, interface compatibility.
