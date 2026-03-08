package main

import (
	"fmt"
	"runtime"
)

// getSmallSliceLeak creates a large slice but returns a small sub-slice.
// The underlying array of the large slice is NOT garbage collected because
// the returned small slice still references it.
//
//go:noinline
func getSmallSliceLeak() []int {
	// Allocate a large slice (e.g., 10 million ints ~ 80MB)
	largeSlice := make([]int, 10_000_000)
	for i := 0; i < len(largeSlice); i++ {
		largeSlice[i] = i
	}

	// Return a tiny slice of the first 5 elements.
	// PROBLEM: The backing array for these 5 elements is the same 1 million element array.
	// Even though 'largeSlice' goes out of scope, the runtime sees that the returned
	// slice points to the same memory block, so the WHOLE block (8MB) stays in memory.
	return largeSlice[:5]
}

// getSmallSliceFixed uses 'copy' to ensure we only hold memory for what we need.
//
//go:noinline
func getSmallSliceFixed() []int {
	largeSlice := make([]int, 10_000_000)
	for i := 0; i < len(largeSlice); i++ {
		largeSlice[i] = i
	}

	// FIX: Create a new independent slice of the desired size.
	smallSlice := make([]int, 5)

	// Copy data from the large array to the new small array.
	copy(smallSlice, largeSlice[:5])

	// Now 'largeSlice' can be garbage collected when this function exits,
	// because 'smallSlice' has its own separate backing array.
	return smallSlice
}

func printMemUsage(msg string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// Alloc = bytes of allocated heap objects.
	fmt.Printf("%s: Alloc = %v KB (%v MiB)\n", msg, m.Alloc/1024, m.Alloc/1024/1024)
	// Note: In some environments/Go versions, compiler optimizations might mask the leak
	// in this simple example, but the structural leak is real.
}

func main() {
	fmt.Println("--- Demonstrating Slice Memory Leak ---")

	// Force GC to start clean
	runtime.GC()
	printMemUsage("Initial Memory")

	// Scenario 1: The Leak
	fmt.Println("\n[Calling Leaky Function...]")
	leak := getSmallSliceLeak()
	// We need to keep 'leak' alive so the compiler doesn't optimize it away
	fmt.Println("Leaky slice length:", len(leak))

	runtime.GC() // Try to clean up
	printMemUsage("After Leaky Function (Large array still in memory)")

	// Clear reference to allow GC
	leak = nil
	runtime.GC()
	printMemUsage("After clearing leak (Memory reclaimed)")

	fmt.Println("\n---------------------------------------")

	// Scenario 2: The Fix
	fmt.Println("[Calling Fixed Function...]")
	fixed := getSmallSliceFixed()
	fmt.Println("Fixed slice length:", len(fixed))

	runtime.GC() // Try to clean up
	printMemUsage("After Fixed Function (Large array collected immediately)")

	// Keep fixed alive just for the sake of the example
	_ = fixed
}
