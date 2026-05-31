// generics_advanced.go
// Generics beyond the toy example: constraints, inference, when NOT to use,
// and the stdlib packages that ship with them (`slices`, `maps`, `cmp`).
package main

import (
	"cmp"
	"fmt"
	"slices"
)

// ---------------------------------------------------------------------------
// 1. Constraints — what types T is allowed to be.
// ---------------------------------------------------------------------------

// `any` — no constraint (alias for interface{}).
func First[T any](s []T) (T, bool) {
	var zero T
	if len(s) == 0 {
		return zero, false
	}
	return s[0], true
}

// `comparable` — types usable with == and !=
// (no slices, maps, funcs; structs are comparable iff all fields are).
func Contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// Custom constraint via interface union.
// The `~` means "underlying type" — accepts named types too.
//   type MyInt int   // matches ~int, not int alone
type Ordered interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Use the stdlib `cmp.Ordered` (Go 1.21+) instead of rolling your own.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// ---------------------------------------------------------------------------
// 2. Type inference — when can you omit the type argument?
// ---------------------------------------------------------------------------

// The compiler infers T from arguments. These are equivalent:
//   Max(3, 5)        // inferred  T=int
//   Max[int](3, 5)   // explicit

// Inference works for function arguments but NOT for return-type-only generics.
// You must specify T explicitly when T appears only in the return:
//   func Zero[T any]() T { var z T; return z }
//   v := Zero[int]()  // must say [int]

// ---------------------------------------------------------------------------
// 3. Generic types and methods.
// Methods on generic types CANNOT introduce new type parameters.
// ---------------------------------------------------------------------------

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}
func (s *Stack[T]) Len() int { return len(s.items) }

// ---------------------------------------------------------------------------
// 4. The stdlib generic packages — know these for interviews.
// ---------------------------------------------------------------------------

func stdlibTour() {
	xs := []int{3, 1, 4, 1, 5, 9, 2, 6}

	slices.Sort(xs) // in-place, ascending
	idx, found := slices.BinarySearch(xs, 4)
	fmt.Println("binary search 4:", idx, found)

	fmt.Println("contains 9:", slices.Contains(xs, 9))
	fmt.Println("max:", slices.Max(xs))
	fmt.Println("equal:", slices.Equal(xs, []int{1, 1, 2, 3, 4, 5, 6, 9}))

	// slices.SortFunc with a comparator returning -1, 0, +1.
	type user struct {
		name string
		age  int
	}
	users := []user{{"alice", 30}, {"bob", 25}, {"carol", 35}}
	slices.SortFunc(users, func(a, b user) int { return cmp.Compare(a.age, b.age) })
	fmt.Println("by age:", users)
}

func main() {
	v, _ := First([]string{"a", "b"})
	fmt.Println("first:", v)

	fmt.Println("contains:", Contains([]int{1, 2, 3}, 2))
	fmt.Println("max:", Max(3.5, 2.1))
	fmt.Println("min:", Min("apple", "banana"))

	s := &Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	top, _ := s.Pop()
	fmt.Println("popped:", top, "len:", s.Len())

	stdlibTour()
}

/*
WHEN TO USE GENERICS
  - Container types (Stack, Set, LRU cache) where the element type varies
    but the operations don't.
  - Algorithms over orderable/comparable types (Min, Max, Contains, Sort).
  - Type-safe collections that previously needed `any` + assertions.

WHEN NOT TO USE GENERICS
  - When a single concrete type is enough — generics add cognitive cost.
  - When an interface fits naturally (e.g., io.Reader). Generics are not a
    replacement for polymorphism via interfaces.
  - When you'd need to type-switch on T inside the body. That's a smell —
    you're really writing per-type code; just write the types out.
  - Hot paths where escape-analysis matters. Generic code can box more.

KEY RULES (interview probes)
  - Methods on generic types can't add new type parameters.
  - `comparable` allows == and !=. Doesn't include slice/map/func types
    or structs containing them.
  - `~T` means "any type whose underlying type is T" — needed to accept
    named types like `type Celsius float64`.
  - `cmp.Ordered` (Go 1.21+) is the stdlib constraint for orderable types.
    Prefer it over a hand-rolled union.
  - No method-level type parameters: `func (s *S) M[T any](t T)` won't compile.
  - Type inference is one-way: from arguments to type params, not from
    return-type usage.

STDLIB GENERIC PACKAGES (memorize)
  slices  — Sort, SortFunc, BinarySearch, Contains, Index, Max, Min, Equal,
            Reverse, Clone, Delete, Insert, Compact, Concat
  maps    — Keys, Values, Equal, Clone, Copy, DeleteFunc
  cmp     — Ordered, Compare, Less, Or
*/
