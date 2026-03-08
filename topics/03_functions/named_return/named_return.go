// named_return.go
// Demonstrates named return values in Go.
package main

import "fmt"

// Named return values
func divide(a, b float64) (result float64, err error) {
	if b == 0 {
		err = fmt.Errorf("division by zero")
		return
	}
	result = a / b
	return
}

func main() {
	res, err := divide(10, 2)
	fmt.Println("10 / 2 =", res, "error:", err)
	res, err = divide(10, 0)
	fmt.Println("10 / 0 =", res, "error:", err)
}
// Documentation:
// - Named return values are initialized to zero values.
// - Use bare return to return current values.
// - Edge cases: error handling, readability concerns.
