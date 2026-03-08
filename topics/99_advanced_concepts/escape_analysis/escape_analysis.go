// escape_analysis.go
// Demonstrates how values escape from stack to heap.
// Run with: go build -gcflags="-m" escape_analysis.go
// Lines annotated with "escapes to heap" in the -m output are key.
package main

import "fmt"

// Case 1: returning a pointer forces the local var to heap
func newValue() *int {
	x := 42 // x escapes to heap because its address is returned
	return &x
}

// Case 2: closures capture variables — captured variables escape to heap
func makeAdder(base int) func(int) int {
	// base escapes to heap because the closure outlives this stack frame
	return func(n int) int {
		return base + n
	}
}

// Case 3: interface boxing — concrete value is copied onto the heap
// when assigned to an interface value
type Animal interface {
	Sound() string
}

type Dog struct{ Name string }

func (d Dog) Sound() string { return "woof" }

func boxIntoInterface(d Dog) Animal {
	return d // d's value is copied; the copy escapes to heap
}

// Case 4: large struct — Go may decide it's cheaper to put large objects on heap
func largeOnStack() {
	var buf [1024]byte // small enough to stay on stack
	buf[0] = 1
	_ = buf
}

func main() {
	p := newValue()
	fmt.Println("heap-allocated int:", *p)

	add5 := makeAdder(5)
	fmt.Println("closure result:", add5(10))

	a := boxIntoInterface(Dog{Name: "Rex"})
	fmt.Println("interface value:", a.Sound())

	largeOnStack()
	fmt.Println("large array stayed on stack (check -gcflags output to verify)")
}
