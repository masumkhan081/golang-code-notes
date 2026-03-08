// subtests_test.go
// Run with: go test -v -run TestDivSubtest ./10_testing/
// NOT runnable with go run — this is a test file.
package main

import (
	"fmt"
	"testing"
)

func Multiply(a, b int) int {
	return a * b
}

func TestMultiply(t *testing.T) {
	// Subtests allow you to run specific tests using -run TestMultiply/Name

	t.Run("positive numbers", func(t *testing.T) {
		if got := Multiply(2, 3); got != 6 {
			t.Errorf("Multiply(2, 3) = %d; want 6", got)
		}
	})

	t.Run("negative numbers", func(t *testing.T) {
		if got := Multiply(-2, 3); got != -6 {
			t.Errorf("Multiply(-2, 3) = %d; want -6", got)
		}
	})

	t.Run("zeros", func(t *testing.T) {
		if got := Multiply(5, 0); got != 0 {
			t.Errorf("Multiply(5, 0) = %d; want 0", got)
		}
	})

	// Table-driven with subtests
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"2x2", 2, 2, 4},
		{"3x3", 3, 3, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Multiply(tt.a, tt.b); got != tt.want {
				t.Errorf("Multiply(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func ExampleMultiply() {
	fmt.Println(Multiply(2, 5))
	// Output: 10
}
