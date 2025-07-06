// interfaces_advanced.go
// Demonstrates Go interfaces, type assertion, type switch, and empty interface usage.
package main

import "fmt"

// Speaker interface defines the behavior of speaking
type Speaker interface {
	Speak() string
}

// Dog struct
type Dog struct {
	Name string
}

// Speak method for Dog, implicitly implements Speaker
func (d Dog) Speak() string {
	return "Woof! My name is " + d.Name
}

// Cat struct
type Cat struct {
	Name string
}

// Speak method for Cat, implicitly implements Speaker
func (c Cat) Speak() string {
	return "Meow! My name is " + c.Name
}

// Describe function takes a Speaker interface
// It can work with any type that implements the Speaker interface
func Describe(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	dog := Dog{Name: "Buddy"}
	cat := Cat{Name: "Whiskers"}

	Describe(dog) // Output: Woof! My name is Buddy
	Describe(cat) // Output: Meow! My name is Whiskers

	// A slice of interfaces
	animals := []Speaker{dog, cat}
	for _, animal := range animals {
		fmt.Println(animal.Speak())
	}

	// Type Assertion (with error checking)
	var s Speaker = dog // Assign a concrete type to an interface variable
	if d, ok := s.(Dog); ok {
		fmt.Printf("It's a dog named %s\n", d.Name)
	} else {
		fmt.Println("Not a dog.")
	}

	// Type Switch
	checkAnimal := func(s Speaker) {
		switch v := s.(type) {
		case Dog:
			fmt.Printf("Detected a Dog: %s\n", v.Name)
		case Cat:
			fmt.Printf("Detected a Cat: %s\n", v.Name)
		default:
			fmt.Printf("Unknown speaker type: %T\n", v)
		}
	}

	checkAnimal(dog)
	checkAnimal(cat)
	// Example with an empty interface (not related to Speaker, just for demonstration)
	var i interface{} = 10
	var j interface{} = "hello"
	fmt.Println("Empty interface value i:", i)
	fmt.Println("Empty interface value j:", j)
}
