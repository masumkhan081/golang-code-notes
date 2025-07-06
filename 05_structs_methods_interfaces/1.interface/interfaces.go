// interfaces.go
// Demonstrates interfaces, implementation, and edge cases in Go.
package main

import "fmt"

// Describer interface
type Describer interface {
	Describe() string
}

// User struct implements Describer
type User struct {
	Name string
}

func (u User) Describe() string {
	return "User: " + u.Name
}

// Product struct implements Describer
type Product struct {
	ID int
}

func (p Product) Describe() string {
	return fmt.Sprintf("Product ID: %d", p.ID)
}

func main() {
	var d Describer
	d = User{Name: "Alice"}
	fmt.Println(d.Describe())
	d = Product{ID: 101}
	fmt.Println(d.Describe())

	// Edge case: interface nil value
	var n Describer
	fmt.Println("Nil interface:", n)
}
// Documentation:
// - Interfaces are satisfied implicitly.
// - Any type that implements the methods is compatible.
// - Edge cases: nil interfaces, empty interfaces, type assertions.
