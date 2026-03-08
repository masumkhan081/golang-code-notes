// error_type.go
// Demonstrates basic error handling and the error type in Go.
package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 2 =", result)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
// Documentation:
// - The built-in error type is an interface.
// - Return error as the last value for idiomatic Go error handling.
// - Edge cases: nil error, error propagation.
