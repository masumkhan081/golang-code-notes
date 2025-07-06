// main_test.go
// Demonstrates unit test, table-driven test, and benchmark in Go.
package main

import (
	"testing"
)

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{1, 2, 3},
		{2, 2, 4},
		{0, 0, 0},
	}
	for _, tt := range tests {
		got := Add(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(1, 2)
	}
}
// Documentation:
// - Use *_test.go files for tests.
// - Table-driven tests are idiomatic for multiple cases.
// - Benchmarks use BenchmarkXxx naming.
// - Edge cases: test failures, coverage, setup/teardown.
