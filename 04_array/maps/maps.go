// maps.go
// Demonstrates maps in Go, including edge cases.
package main

import "fmt"

func main() {
	// Declare and initialize a map
	m := map[string]int{"one": 1, "two": 2}
	fmt.Println("m:", m)

	// Add and update
	m["three"] = 3
	m["one"] = 11
	fmt.Println("m after add/update:", m)

	// Access and check existence
	val, ok := m["two"]
	fmt.Println("m[\"two\"]:", val, "exists:", ok)

	// Delete
	delete(m, "two")
	fmt.Println("m after delete:", m)

	// Edge case: non-existent key
	v := m["absent"]
	fmt.Println("m[\"absent\"] (zero value):", v)

	// Edge case: nil map
	var nilMap map[string]int
	// nilMap["fail"] = 1 // ERROR: assignment to entry in nil map
	fmt.Println("nilMap:", nilMap)
}
// Documentation:
// - Maps are reference types, keys must be comparable.
// - Accessing a non-existent key returns zero value.
// - Edge cases: assignment to nil map, iteration order is random.
