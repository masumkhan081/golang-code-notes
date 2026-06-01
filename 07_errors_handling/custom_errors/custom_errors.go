// custom_errors.go
// Demonstrates custom error types in Go.
package main

import (
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return ValidationError{Field: "Age", Message: "must be non-negative"}
	}
	return nil
}

func main() {
	err := validateAge(-5)
	if err != nil {
		fmt.Println("Validation error:", err)
		// Type assertion to access custom fields
		if vErr, ok := err.(ValidationError); ok {
			fmt.Println("Field:", vErr.Field, "Message:", vErr.Message)
		}
	}
}
// Documentation:
// - Custom error types allow for richer error information.
// - Implement the Error() string method.
// - Edge cases: error type assertions, nil custom errors.
