// arrays.go
// Demonstrates arrays in Go, including edge cases.
package main

import "fmt"

func main() {
	// Declare and initialize an array
	var arr1 [3]int = [3]int{1, 2, 3}
	fmt.Println("arr1:", arr1)

	// Short declaration
	arr2 := [5]string{"a", "b", "c", "d", "e"}
	fmt.Println("arr2:", arr2)

	// Array of zero values
	var arr3 [2]float64
	fmt.Println("arr3 (zero values):", arr3)

	// Array length is part of type
	// var arr4 [4]int = arr1 // ERROR: cannot use arr1 (type [3]int) as type [4]int

	// Edge case: array copy is by value
	arrCopy := arr1
	arrCopy[0] = 99
	fmt.Println("arr1 after copy-modify:", arr1)
	fmt.Println("arrCopy:", arrCopy)
}
// Documentation:
// - Arrays have fixed length and are value types.
// - Assignment and passing copies the entire array.
// - Edge cases: type mismatch on length, memory cost for large arrays.
