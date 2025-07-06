// data_types.go
// Examples of Go data types
// Demonstrates basic, composite, and edge-case data types in Go.
package main

import (
	"fmt"
	"math"
)

type MyInt int

type MyStruct struct {
	Name string
	Age  int
}

func main() {
	var a int = 42
	var b int8 = -128
	var c uint16 = 65535
	fmt.Println("int:", a, "int8:", b, "uint16:", c)

	var f float64 = math.Pi
	fmt.Printf("float64: %.5f\n", f)

	var isGoFun bool = true
	fmt.Println("bool:", isGoFun)

	var s string = "GoLang"
	fmt.Println("string:", s)

	var r rune = '♥'
	fmt.Printf("rune: %c (%U)\n", r, r)

	var by byte = 255
	fmt.Println("byte:", by)

	arr := [3]int{1, 2, 3}
	fmt.Println("array:", arr)

	slc := []string{"a", "b", "c"}
	fmt.Println("slice:", slc)

	m := map[string]int{"one": 1, "two": 2}
	fmt.Println("map:", m)

	person := MyStruct{Name: "Alice", Age: 30}
	fmt.Println("struct:", person)

	p := &a
	fmt.Println("pointer to a:", *p)

	var zeroInt int
	var zeroFloat float64
	var zeroBool bool
	var zeroString string
	fmt.Printf("zero values: int=%d, float=%f, bool=%v, string='%s'\n", zeroInt, zeroFloat, zeroBool, zeroString)

	var nilSlice []int
	var nilMap map[string]int
	fmt.Println("nil slice:", nilSlice, "len:", len(nilSlice))
	fmt.Println("nil map:", nilMap)

	var my MyInt = 100
	fmt.Println("custom type MyInt:", my)
}
// Documentation:
// - Go supports strong, static typing.
// - Zero values are important for understanding default initialization.
// - Nil is only valid for reference types (slice, map, pointer, channel, function, interface).
// - Use type aliasing and custom types for clarity and safety.
// - Edge cases: overflow, underflow, nil dereference (not shown here).
