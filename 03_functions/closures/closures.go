// closures.go
// Demonstrates closures in Go.
package main

import "fmt"

func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	counter := makeCounter()
	fmt.Println(counter()) // 1
	fmt.Println(counter()) // 2
	fmt.Println(counter()) // 3

	// Edge case: closure capturing loop variable
	funcs := []func(){ }
	for i := 0; i < 3; i++ {
		val := i // capture current value
		funcs = append(funcs, func() { fmt.Println("val:", val) })
	}
	for _, f := range funcs {
		f()
	}
}
// Documentation:
// - Closures capture variables from their surrounding scope.
// - Edge cases: capturing loop variables, memory leaks if closures outlive their context.
