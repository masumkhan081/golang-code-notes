package main

import "fmt"

// MyError is a custom error struct
type MyError struct {
	Code int
}

// Error implements the error interface
func (e *MyError) Error() string {
	return fmt.Sprintf("Error Code: %d", e.Code)
}

// returnsNilError returns a *MyError that is nil.
// This is a common trap!
func returnsNilError(fail bool) error {
	var p *MyError = nil
	if fail {
		p = &MyError{Code: 500}
	}
	// WARNING: returning a typed nil pointer (*MyError) as an interface (error)
	// results in a non-nil interface value.
	// The interface will hold: (type=*MyError, value=nil)
	return p
}

// returnsProperNil returns a literal nil interface.
func returnsProperNil(fail bool) error {
	if fail {
		return &MyError{Code: 500}
	}
	// CORRECT: Explicitly return nil for the interface
	// The interface will hold: (type=nil, value=nil)
	return nil
}

func main() {
	fmt.Println("--- The Nil Interface Trap ---")

	// Case 1: The Trap
	err := returnsNilError(false)

	fmt.Printf("Case 1: returnsNilError(false)\n")
	fmt.Printf("  Value: %v\n", err)

	// The check 'err != nil' is TRUE, even though the underlying value is nil!
	if err != nil {
		fmt.Println("  Result: err != nil (TRAP! We expected nil)")
		fmt.Printf("  Type of err: %T\n", err)
	} else {
		fmt.Println("  Result: err == nil (Correct)")
	}

	fmt.Println("\n------------------------------")

	// Case 2: The Fix
	err2 := returnsProperNil(false)

	fmt.Printf("Case 2: returnsProperNil(false)\n")
	fmt.Printf("  Value: %v\n", err2)

	if err2 != nil {
		fmt.Println("  Result: err2 != nil (Unexpected)")
	} else {
		fmt.Println("  Result: err2 == nil (Correct! This is what we want)")
	}

	fmt.Println("\n--- Explanation ---")
	fmt.Println("An interface in Go is a tuple: (type, value).")
	fmt.Println("A nil interface is (nil, nil).")
	fmt.Println("In Case 1, we returned (*MyError, nil). Since the type is not nil, the interface is not nil.")
	fmt.Println("Always return explicit 'nil' for interfaces if there is no value.")
}
