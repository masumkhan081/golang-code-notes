// interface_internal_representation.go
// Demonstrates dynamic type, dynamic value, nil interface vs typed nil.
package main

import "fmt"

type Animal interface {
	Sound() string
}

type Cat struct{ Name string }

func (c *Cat) Sound() string { return "meow" }

// typed nil trap: a non-nil interface holding a nil pointer
func noCat() Animal {
	var c *Cat // nil pointer of type *Cat
	return c   // interface is NOT nil — it has a dynamic type (*Cat) but nil dynamic value
}

func nilInterface() Animal {
	return nil // truly nil interface: both type and value are nil
}

func inspectInterface(a Animal) {
	if a == nil {
		fmt.Println("  interface is nil (no type, no value)")
		return
	}
	fmt.Printf("  interface is non-nil: dynamic value = %v\n", a)
	// calling a method on a nil *Cat pointer is fine only if the method handles it
}

func main() {
	fmt.Println("=== nil interface ===")
	i1 := nilInterface()
	inspectInterface(i1)

	fmt.Println("\n=== typed nil (common trap) ===")
	i2 := noCat()
	inspectInterface(i2)
	// i2 != nil even though the *Cat inside is nil
	fmt.Println("  i2 == nil?", i2 == nil) // false — this surprises many devs

	fmt.Println("\n=== concrete non-nil value ===")
	i3 := Animal(&Cat{Name: "Whiskers"})
	inspectInterface(i3)
	fmt.Println("  i3 == nil?", i3 == nil)

	fmt.Println("\n=== interface copies concrete value ===")
	// value types are copied into the interface; pointer types store the pointer
	type Point struct{ X, Y int }
	type Shaper interface{ Area() float64 }
	// Point does not implement Shaper here — just demonstrating boxing
	p := &Cat{Name: "Box"}
	var a Animal = p
	p.Name = "Changed"
	// interface holds the pointer, so change is visible through the interface
	fmt.Println("  changed via pointer:", a.Sound())
}
