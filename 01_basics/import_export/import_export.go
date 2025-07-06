// import_export.go
// Examples of import/export in Go
// Demonstrates import statements, aliases, dot imports, and export rules.
package main

import (
	"fmt"
	alias "math/rand"
	. "strings"
)

// Exported identifier (starts with uppercase)
func ExportedFunc() {
	fmt.Println("This function is exported and can be used in other packages.")
}

// unexported identifier (starts with lowercase)
func unexportedFunc() {
	fmt.Println("This function is not exported.")
}

func main() {
	fmt.Println("Standard import: fmt")
	fmt.Println("Random int (alias):", alias.Intn(10))
	fmt.Println("Dot import ToUpper:", ToUpper("go import/export"))

	ExportedFunc()
	unexportedFunc()

	// Edge case: import cycle (will cause compile error if you try to import this package back)
	// Edge case: blank import (use _ \"package\" for side effects)
	// import _ "net/http/pprof"
}
// Documentation:
// - Only identifiers starting with uppercase are exported.
// - Import cycles are not allowed in Go.
// - Aliases and dot imports can be useful but may reduce clarity.
// - Blank imports are used for side effects (e.g., database/sql drivers).
// - Edge cases: import cycles, name collisions, dot import ambiguity.
