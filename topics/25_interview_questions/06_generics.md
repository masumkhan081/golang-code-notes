# Generics (Go 1.18+)

Rarely the centerpiece of a Go interview, but expected basic literacy.
You will be asked: "have you used them, when would you?"

---

### Q1. What problem do generics solve?

Type-safe reuse without `interface{}` + assertions. Pre-generics, a generic
`Contains` had to be written per type or use `any` and lose type safety.
Now: one function, compile-time checked.

```go
func Contains[T comparable](s []T, v T) bool { ... }
```

---

### Q2. What's the difference between `any`, `comparable`, and `cmp.Ordered`?

- `any` — no constraint (alias for `interface{}`).
- `comparable` — supports `==` and `!=`. Excludes slices, maps, funcs, and
  structs containing them. Required for map keys.
- `cmp.Ordered` (Go 1.21+) — supports `<`, `<=`, `>`, `>=`. Integers,
  floats, strings.

---

### Q3. What does `~T` mean in a constraint?

"Any type whose **underlying type** is T." Without it, the constraint
matches only the exact type:

```go
type Celsius float64
type OrderedF interface{ float64 }      // does NOT accept Celsius
type OrderedF2 interface{ ~float64 }    // accepts Celsius
```

You almost always want `~`. The stdlib `cmp.Ordered` uses it.

---

### Q4. Can a method on a generic type add its own type parameters?

**No.** Methods inherit the receiver's type parameters and cannot add new ones:

```go
type Stack[T any] struct { ... }

// COMPILE ERROR — methods can't introduce new type parameters
func (s *Stack[T]) MapTo[U any](f func(T) U) []U { ... }
```

Workaround: make it a free function, not a method.

```go
func MapStack[T, U any](s *Stack[T], f func(T) U) []U { ... }
```

This is the single most common interview gotcha for generics.

---

### Q5. When does Go infer the type argument and when do you have to spell it out?

Inferred from function **arguments**:
```go
Max(3, 5)       // T = int
Max[float64](3, 5) // explicit; converts
```

Required when T appears **only** in the return type:
```go
func Zero[T any]() T { var z T; return z }
v := Zero[int]() // must specify
```

---

### Q6. Why are generics not always faster?

Go uses **GCShape stenciling** — one compiled body per "shape" of T
(roughly: types with identical layout share code). Pointer-typed instantiations
share one body; value types may also share. This trades binary size for
compile-time cost, but:
- The body uses a dictionary to look up methods at runtime — slower than
  hand-written code for the concrete type.
- Inlining is harder for generic functions.

For hot paths, benchmark before assuming generics win. Often `interface{}`
or a hand-specialized function is comparable or faster.

---

### Q7. When should you NOT use generics?

- One concrete type is enough — generics add reading cost.
- An interface fits naturally (`io.Reader`, `error`) — that's polymorphism.
- You'd need to type-switch on T inside the body — you're really writing
  per-type code; just write multiple functions.
- Tight, allocation-sensitive paths until you benchmark.

---

### Q8. What stdlib packages use generics?

- **`slices`** — `Sort`, `SortFunc`, `BinarySearch`, `Contains`, `Index`,
  `Max`, `Min`, `Equal`, `Reverse`, `Clone`, `Delete`, `Insert`, `Compact`.
- **`maps`** — `Keys`, `Values`, `Equal`, `Clone`, `Copy`, `DeleteFunc`.
- **`cmp`** — `Ordered`, `Compare`, `Less`, `Or`.
- **`sync/atomic`** — `atomic.Pointer[T]` (Go 1.19+).

Knowing these by name is the easy interview win.

---

### Q9. Can you have a generic interface?

Yes:

```go
type Container[T any] interface {
    Add(T)
    Get(int) T
}
```

Useful for abstract data structures. Less common than generic functions/types.

---

### Q10. What's wrong with this constraint?

```go
type Number interface { int | float64 }
func Sum[T Number](xs []T) T { ... }
```

It rejects `type MyInt int`. Fix:

```go
type Number interface { ~int | ~float64 }
```

Always use `~` in numeric/string constraints unless you have a reason not to.
