// constants.go
// Examples of Go constants
// Demonstrates constant declarations, iota, typed/untyped, and edge cases.
package main

import (
	"fmt"
)

const Pi = 3.14159              // untyped constant
const TypedPi float64 = 3.14159 // typed constant

const (
	A = iota // 0
	B        // 1
	C        // 2
)

const (
	FlagNone  = 0
	FlagRead  = 1 << iota // 1
	FlagWrite             // 2
	FlagExec              // 4
)

const (
	MaxUint8  = 1<<8 - 1
	BigNumber = 1e18
)

func main() {
	fmt.Println("Pi:", Pi)
	fmt.Println("TypedPi:", TypedPi)
	fmt.Println("iota constants:", A, B, C)
	fmt.Println("Bitmask flags:", FlagRead, FlagWrite, FlagExec)
	fmt.Println("MaxUint8:", MaxUint8)
	fmt.Println("BigNumber:", BigNumber)

	var f float32 = Pi
	fmt.Println("float32 from untyped const:", f)

	const Sqrt2 = 1.41421356237
	fmt.Println("Sqrt2:", Sqrt2)

	// Compile-time errors: uncommenting below will cause error
	// const Bad = math.Sqrt(2) // ERROR: math.Sqrt(2) is not a constant
}
// Documentation:
// - Constants are evaluated at compile time.
// - iota is useful for enums and bitmasks.
// - Untyped constants are more flexible.
// - Only constant expressions are allowed (no function calls).
// - Edge cases: overflow, type compatibility, iota resets in new const blocks.
