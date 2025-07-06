// slices.go
// Demonstrates slices in Go, including edge cases.
package main

import "fmt"

func main() {
	// Create a slice from an array
	arr := [5]int{10, 20, 30, 40, 50}
	slc := arr[1:4]
	fmt.Println("slice from arr:", slc)

	// Slice literal
	slc2 := []string{"go", "is", "fun"}
	fmt.Println("slc2:", slc2)

	// Append to slice
	slc2 = append(slc2, "!")
	fmt.Println("slc2 after append:", slc2)

	// Length and capacity
	fmt.Println("len:", len(slc2), "cap:", cap(slc2))

	// Edge case: nil slice
	var nilSlice []int
	fmt.Println("nilSlice:", nilSlice, "len:", len(nilSlice))

	// Edge case: reslicing
	sub := slc2[1:]
	fmt.Println("subslice:", sub)
}
// Documentation:
// - Slices are reference types and can grow dynamically.
// - Appending may reallocate and break reference to original array.
// - Edge cases: nil slices, reslicing beyond capacity, sharing underlying array.
