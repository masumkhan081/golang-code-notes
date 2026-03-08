// slice_growth_and_capacity.go
// Demonstrates append reallocation, capacity growth, and shared backing arrays.
package main

import "fmt"

func appendGrowth() {
	var s []int
	prevCap := 0
	for i := 0; i < 20; i++ {
		s = append(s, i)
		if cap(s) != prevCap {
			fmt.Printf("len=%2d cap=%2d  (capacity grew)\n", len(s), cap(s))
			prevCap = cap(s)
		}
	}
}

func sharedBackingArray() {
	original := make([]int, 3, 6) // len=3, cap=6
	original[0], original[1], original[2] = 10, 20, 30

	// a and b share the same backing array as original
	a := original[:2] // [10, 20]
	b := original[1:] // [20, 30]

	fmt.Println("before mutation — a:", a, "b:", b, "original:", original)

	a[1] = 999 // mutates index 1 of the shared backing array

	fmt.Println("after a[1]=999   — a:", a, "b:", b, "original:", original)
	// b[0] and original[1] are also 999 now — same backing array

	// Appending beyond cap creates a NEW backing array: no more sharing
	a = append(a, 1, 2, 3, 4, 5) // exceeds original cap=6 → reallocates
	a[0] = 0
	fmt.Println("after realloc    — a:", a, "original:", original)
	// original is unchanged
}

func sliceHeaderCopy() {
	s := []int{1, 2, 3}
	// Passing slice to function passes the header (ptr + len + cap) by value,
	// so the function can mutate elements but cannot grow the caller's slice.
	double(s)
	fmt.Println("after double:", s) // [2 4 6] — mutations visible
}

func double(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

func main() {
	fmt.Println("=== capacity growth ===")
	appendGrowth()

	fmt.Println("\n=== shared backing array ===")
	sharedBackingArray()

	fmt.Println("\n=== slice header copy ===")
	sliceHeaderCopy()
}
