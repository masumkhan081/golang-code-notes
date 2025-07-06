// structs.go
// Demonstrates struct declaration, initialization, pointers, and zero values in Go.
package main

import "fmt"

// Person struct defines a person with Name and Age
type Person struct {
	Name string
	Age  int
}

func main() {
	// 1. Declare and initialize a struct
	var p1 Person
	p1.Name = "Alice"
	p1.Age = 30
	fmt.Println("Person 1:", p1) // Output: Person 1: {Alice 30}

	// 2. Initialize using struct literal (named fields - preferred)
	p2 := Person{Name: "Bob", Age: 25}
	fmt.Println("Person 2:", p2) // Output: Person 2: {Bob 25}

	// 3. Initialize using struct literal (ordered fields - less readable for many fields)
	p3 := Person{"Charlie", 35}
	fmt.Println("Person 3:", p3) // Output: Person 3: {Charlie 35}

	// 4. Pointer to a struct
	p4 := &Person{Name: "David", Age: 40}
	fmt.Println("Person 4 (pointer):", *p4) // Output: Person 4 (pointer): {David 40}
	fmt.Println("Person 4 Name (via pointer):", p4.Name) // Go automatically dereferences

	// 5. Zero values of a struct
	var p5 Person
	fmt.Println("Person 5 (zero values):", p5) // Output: Person 5 (zero values): { 0}
}
