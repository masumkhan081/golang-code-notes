// interface_idioms.go
// Go interface design rules that come up in every interview.
//
// 1. any vs interface{}
// 2. Accept interfaces, return structs
// 3. Keep interfaces small; define them at the consumer
// 4. Pointer vs value receivers and method set rules
//
// For the typed-nil gotcha (most-asked interface trap), see
// 99_advanced_concepts/nil_interface_trap/nil_interface_trap.go
package main

import (
	"fmt"
	"io"
	"strings"
)

// ---------------------------------------------------------------------------
// 1. any IS interface{}
// `any` is a Go 1.18+ alias for `interface{}`. Identical at runtime.
// Use `any` in new code — it's shorter and signals intent ("any type"
// vs "empty interface", which sounds technical).
// ---------------------------------------------------------------------------

func printAny(v any)          { fmt.Println(v) }
func printEmpty(v interface{}) { fmt.Println(v) } // same thing

// ---------------------------------------------------------------------------
// 2. "Accept interfaces, return structs"
// Functions take the smallest interface they need (flexibility for callers)
// and return concrete types (callers see all methods, no hidden behavior).
// ---------------------------------------------------------------------------

// GOOD: takes io.Reader (any source), returns concrete *strings.Builder.
func collect(r io.Reader) *strings.Builder {
	var b strings.Builder
	buf := make([]byte, 64)
	for {
		n, err := r.Read(buf)
		b.Write(buf[:n])
		if err != nil {
			return &b
		}
	}
}

// BAD: takes a concrete *strings.Reader — callers can't pass a file, network, etc.
// BAD: returns an interface — callers lose access to type-specific methods.

// ---------------------------------------------------------------------------
// 3. Small interfaces, defined at the consumer.
// io.Reader is one method. io.Writer is one method. They compose into
// io.ReadWriter when needed. This is the Go style.
//
// Interfaces should be DECLARED by the package that USES them,
// not by the package that implements them. The implementer doesn't need
// to know which interfaces it satisfies — Go matches structurally.
// ---------------------------------------------------------------------------

// Consumer-side interface. The "store" package implementing this never
// imports this file and never says "implements UserFinder".
type UserFinder interface {
	FindUser(id string) (string, bool)
}

func Greet(uf UserFinder, id string) string {
	if name, ok := uf.FindUser(id); ok {
		return "hello, " + name
	}
	return "unknown user"
}

// A concrete implementer — note no interface mention.
type memStore map[string]string

func (m memStore) FindUser(id string) (string, bool) {
	v, ok := m[id]
	return v, ok
}

// ---------------------------------------------------------------------------
// 4. Method sets — what satisfies what?
//
// Given:  func (t T)  M1()      // value receiver
//         func (t *T) M2()      // pointer receiver
//
// Method set of  T:  { M1 }
// Method set of *T:  { M1, M2 }
//
// A *T satisfies any interface T does, plus interfaces requiring M2.
// A  T satisfies only interfaces requiring M1.
//
// You CAN call M2 on a T variable (Go auto-addresses) — but only if T
// is addressable. T values inside interface boxes are NOT addressable.
// ---------------------------------------------------------------------------

type Counter struct{ n int }

func (c Counter) Value() int { return c.n }
func (c *Counter) Inc()      { c.n++ }

type Incrementer interface{ Inc() }

func bumpIfPossible(v any) {
	if i, ok := v.(Incrementer); ok {
		i.Inc()
		fmt.Println("bumped")
		return
	}
	fmt.Println("not incrementable")
}

func main() {
	printAny(42)
	printEmpty("same")

	out := collect(strings.NewReader("hello world"))
	fmt.Println("collected:", out.String())

	store := memStore{"u1": "Alice"}
	fmt.Println(Greet(store, "u1"))

	// Method set surprise:
	c := Counter{}
	bumpIfPossible(c)  // not incrementable — Counter (value) lacks Inc
	bumpIfPossible(&c) // bumped         — *Counter has Inc
	fmt.Println("c.n =", c.n)
}

/*
RULES OF THUMB
  - Define interfaces in the package that consumes them. Implementer stays free.
  - Smaller is better. One-method interfaces compose well.
  - Use `any`, not `interface{}`, in new code.
  - Return concrete types unless you have a clear reason not to.
  - Use *T receivers consistently per type — don't mix value and pointer
    receivers on the same type unless you have a reason.

ASSERTION FORMS
  v, ok := i.(T)      // comma-ok: safe, never panics
  v := i.(T)          // panics if i is not T
  switch v := i.(type) // type switch — preferred over chains of asserts

COMMON TRAP — TYPED NIL
  See 99_advanced_concepts/nil_interface_trap. In one line:
    var p *MyErr = nil
    var e error = p       // e != nil — interface has a type slot even if value is nil
*/
