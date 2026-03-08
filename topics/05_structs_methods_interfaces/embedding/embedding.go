// embedding.go
// Demonstrates struct embedding and method promotion in Go.
package main

import "fmt"

// Base struct
type Animal struct {
	Name string
}

func (a Animal) Speak() {
	fmt.Println(a.Name, "makes a sound")
}

// Dog embeds Animal
type Dog struct {
	Animal
	Breed string
}

func (d Dog) Bark() {
	fmt.Println(d.Name, "says woof!")
}

func main() {
	dog := Dog{Animal: Animal{Name: "Buddy"}, Breed: "Beagle"}
	dog.Speak() // Promoted method
	dog.Bark()
	fmt.Println("Dog struct:", dog)
}
// Documentation:
// - Embedding promotes fields and methods to the outer struct.
// - Outer struct can override embedded methods.
// - Edge cases: field/method name conflicts, deep embedding.
