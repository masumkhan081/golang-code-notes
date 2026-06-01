// wrapping.go
// Demonstrates error wrapping and unwrapping in Go.
package main

import (
	"errors"
	"fmt"
)

func readFile(filename string) error {
	return fmt.Errorf("failed to read %s: %w", filename, errors.New("file not found"))
}

func main() {
	err := readFile("test.txt")
	if err != nil {
		fmt.Println("Wrapped error:", err)
		// Unwrap to get the original error
		unwrapped := errors.Unwrap(err)
		fmt.Println("Unwrapped error:", unwrapped)
	}
}
// Documentation:
// - Use %w with fmt.Errorf to wrap errors.
// - Use errors.Unwrap or errors.Is/As to inspect wrapped errors.
// - Edge cases: multiple wraps, error chain traversal.
