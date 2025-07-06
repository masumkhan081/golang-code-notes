// variables.go
// Examples of Go variables
// Demonstrates declaration, initialization, shadowing, and edge cases.
package main

import "fmt"

var globalVar = "I am global"

func main() {
	var a int = 10
	var b = 20
	c := 30
	fmt.Println("a:", a, "b:", b, "c:", c)

	var x, y, z int = 1, 2, 3
	fmt.Println("x, y, z:", x, y, z)

	var zeroInt int
	fmt.Println("zeroInt:", zeroInt)

	var (
		str string = "block"
		flag bool   = true
	)
	fmt.Println("str:", str, "flag:", flag)

	fmt.Println("globalVar:", globalVar)
	globalVar := "I am shadowed locally"
	fmt.Println("shadowed globalVar:", globalVar)

	// Edge case: unused variables (will cause compile error)
	// var unused int // ERROR: unused variable

	// Edge case: redeclaration in same scope (will cause compile error)
	// var a int = 100 // ERROR: a redeclared in this block
}
// Documentation:
// - Variables must be used, or compilation fails.
// - Short declaration (:=) is not allowed at package level.
// - Shadowing can lead to subtle bugs.
// - Zero values are important for initialization.
// - Edge cases: unused variables, redeclaration, type inference.
