# Interfaces

---

### Q1. How does interface satisfaction work in Go?

**Structurally and implicitly.** A type satisfies an interface if it has
all the required methods — no `implements` keyword, no registration. The
implementer doesn't even need to import the interface.

Consequence: interfaces should be **declared by the consumer**, not the
producer. The smaller the interface, the easier it is to satisfy.

---

### Q2. What's the difference between `interface{}` and `any`?

None at runtime — `any` is a Go 1.18+ alias for `interface{}`. Use `any` in
new code; it reads better.

---

### Q3. The typed-nil trap. What does this print?

```go
type MyErr struct{}
func (*MyErr) Error() string { return "boom" }

func mayFail() error {
    var p *MyErr = nil
    return p
}

func main() {
    err := mayFail()
    fmt.Println(err == nil) // false
}
```

`err == nil` is **false** even though `p` is nil. An interface value is a
**(type, value)** pair. The type slot is `*MyErr`, value is nil — that's
not the same as both slots being nil.

Fix: return `nil` explicitly, never a typed nil pointer:
```go
if !failed {
    return nil
}
return &MyErr{}
```

Code: `24_advanced_concepts/nil_interface_trap/`.

---

### Q4. What's the difference between pointer and value receivers for interface satisfaction?

Given:
```go
func (t T) M1()  // value receiver
func (t *T) M2() // pointer receiver
```

- Method set of `T`  = `{M1}`
- Method set of `*T` = `{M1, M2}`

So `*T` satisfies any interface `T` does, plus interfaces requiring `M2`.
`T` does **not** satisfy interfaces requiring `M2`.

You CAN call `M2` on a `T` variable when it's addressable (Go auto-addresses).
But `T` values **inside an interface** are not addressable — that's where
this rule bites.

---

### Q5. What does "accept interfaces, return structs" mean?

- **Accept interfaces** — take the smallest interface your function needs (e.g., `io.Reader`, not `*os.File`). Callers can pass any matching type.
- **Return structs** — return concrete types so callers can use all methods and the compiler can inline.

Returning an interface hides behavior and prevents the caller from using
type-specific functionality.

---

### Q6. When should you define an interface?

- **At the consumer side**, when you need to substitute behavior (testing, multiple backends).
- When you have **two or more concrete implementations** that need to be interchangeable.

Avoid:
- Defining a one-impl interface "in case" you add more later — YAGNI, and Go interfaces are structural so you can add one when needed without changing the implementer.
- Putting interfaces in the same package as the implementation when only one impl exists.

---

### Q7. Empty interface — when is it appropriate?

Mostly never in new code. With generics (1.18+) you have type-safe options.

Legitimate uses today:
- `fmt.Println(...any)` — variadic of unknown types
- JSON `map[string]any` for unstructured payloads
- Heterogeneous containers (rare)

---

### Q8. What's a type assertion and how does it fail safely?

```go
v, ok := i.(MyType)   // comma-ok form — safe; ok=false if mismatch
v := i.(MyType)       // panics on mismatch
```

For multiple types, prefer a type switch:
```go
switch v := i.(type) {
case int:    // v is int here
case string: // v is string here
default:
}
```

---

### Q9. What's inside an interface value (iface/eface)?

Two words:
- `eface` (empty interface) — `(type *_type, data unsafe.Pointer)`
- `iface` (non-empty)        — `(itab *itab, data unsafe.Pointer)` where `itab` caches the method table for the (interface, concrete type) pair.

Implications:
- Boxing a value in an interface allocates if it doesn't fit in a word.
- Comparing interfaces compares both type and value — and panics if the value type isn't comparable (e.g., contains a slice/map/func).

Code: `24_advanced_concepts/interface_internal_representation/`.

---

### Q10. Why is `io.Reader` such a celebrated interface?

One method: `Read(p []byte) (n int, err error)`. Composes with everything
— files, network sockets, HTTP bodies, gzip streams, hash writers,
test fixtures. Demonstrates the Go style:
- single-method interface
- defined at the consumer (the `io` package)
- composable (`io.ReadWriter`, `io.ReadCloser`)
- accepts wide range of types without those types knowing about it
