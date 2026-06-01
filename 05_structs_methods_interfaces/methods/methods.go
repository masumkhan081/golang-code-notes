// methods.go
// Demonstrates methods on structs, pointer/value receivers, and edge cases in Go.
package main

import "fmt"

// Rectangle struct
type Rectangle struct {
	Width, Height float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Method with pointer receiver (can modify the struct)
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func main() {
	rect := Rectangle{Width: 3, Height: 4}
	fmt.Println("Area:", rect.Area())

	rect.Scale(2)
	fmt.Println("Scaled Rectangle:", rect)
	fmt.Println("Scaled Area:", rect.Area())
}
// Documentation:
// - Methods can have value or pointer receivers.
// - Pointer receivers can modify the struct.
// - Edge cases: method sets, nil receivers, method promotion with embedding.
