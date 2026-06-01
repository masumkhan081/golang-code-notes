// switch.go
// Demonstrates Go's switch statement, including edge cases.
package main

import "fmt"

func main() {
	// Basic switch
	d := 2
	switch d {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	}

	// Switch with multiple expressions
	switch d {
	case 1, 3, 5:
		fmt.Println("odd")
	case 2, 4, 6:
		fmt.Println("even")
	}

	// Type switch
	var x interface{} = 3.14
	switch v := x.(type) {
	case int:
		fmt.Println("int", v)
	case float64:
		fmt.Println("float64", v)
	default:
		fmt.Println("unknown type")
	}

	// Edge case: switch without condition
	switch {
	case d < 0:
		fmt.Println("negative")
	case d == 0:
		fmt.Println("zero")
	default:
		fmt.Println("positive")
	}
}
// Documentation:
// - switch can be value, type, or condition-based.
// - Multiple cases per line allowed.
// - Edge cases: fallthrough, empty switch, type switch.
