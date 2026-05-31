# Misc Go traps and "gotcha" interview questions

Quick-fire questions that test whether you've actually written production Go.

---

### Q1. What does this print?

```go
s := []int{1, 2, 3}
t := s[:1]
t = append(t, 99)
fmt.Println(s)
```

`[1 99 3]` — `t` aliases `s`'s underlying array. `append` had spare
capacity, so it overwrote `s[1]` in place.

Defensive copy when you might mutate:
```go
t := append([]int(nil), s[:1]...)
```

See `99_advanced_concepts/slice_growth_and_capacity/`.

---

### Q2. Why does this leak memory?

```go
big := readBigFile()
small := big[:10]
return small
```

`small` keeps the entire `big` underlying array alive (a slice retains a
pointer to the backing array). To drop the rest:

```go
small := append([]byte(nil), big[:10]...)
```

See `99_advanced_concepts/slice_memory_leak/`.

---

### Q3. What does this print?

```go
m := map[string]int{"a": 1, "b": 2}
for k, v := range m {
    fmt.Println(k, v)
}
```

The iteration order is **randomized per-run** by design — prevents code from accidentally depending on order. Sort keys explicitly if you need stability.

---

### Q4. Why does this fail?

```go
type User struct{ Name string }
var u *User
fmt.Println(u.Name)
```

Panic: nil pointer dereference. But:
```go
type User struct{ Name string }
func (u *User) Hello() string { return "hi" }
var u *User
u.Hello() // does NOT panic — method body just doesn't touch u
```

Method calls on nil receivers are valid Go as long as the body never
dereferences the receiver. Used intentionally for `*MyList.Add()` style APIs.

---

### Q5. What does `defer` capture — value or variable?

**Arguments are evaluated at defer time; receivers/closures capture variables.**

```go
i := 0
defer fmt.Println(i)        // prints 0 — arg captured now
defer func() { fmt.Println(i) }() // prints final i — closure
i = 99
```

Common in benchmarks/timers: `defer track(time.Now())` captures the start
time at the defer line.

---

### Q6. Why doesn't `errors.Is` work here?

```go
if err == sql.ErrNoRows { ... }   // brittle
if errors.Is(err, sql.ErrNoRows) { ... } // correct
```

`==` only works for unwrapped sentinel errors. `errors.Is` walks the wrap
chain (`fmt.Errorf("query: %w", err)`). Same story with `errors.As` for
type assertions through a wrap chain.

---

### Q7. What's a "panic in a deferred function" worth?

A panic during a deferred call **replaces** the original panic — and is
swallowed unless re-panicked. `recover()` only works **inside a deferred
function**, and only catches panics from the same goroutine.

```go
defer func() {
    if r := recover(); r != nil {
        log.Println("recovered:", r)
    }
}()
```

Common interview probe: "can you recover a panic in another goroutine?"
**No.** Each goroutine handles its own. If you spawn goroutines, wrap
each entry point.

---

### Q8. Is `len(channel)` racy?

`len(ch)` and `cap(ch)` are safe to call but the result is **immediately
stale** — by the time you act on it, the value may have changed. Use only
for metrics/diagnostics, never for control flow.

---

### Q9. Map zero-value access vs missing key?

```go
m := map[string]int{"a": 1}
v := m["missing"]  // 0, no panic — returns zero value
v, ok := m["missing"] // 0, false — disambiguate with comma-ok
```

Writing to a nil map panics. Reading from a nil map returns the zero value.

---

### Q10. What's the difference between `new(T)` and `&T{}`?

- `new(T)` — allocates zero value of `T`, returns `*T`.
- `&T{}` — composite literal, returns `*T`. Lets you set fields.

Prefer `&T{}` — it's the same allocation and more flexible. `new` survives mostly for built-in types where you can't write a composite literal (`new(int)`).

---

### Q11. Why does Go discourage `init()`?

- Order across files is alphabetical-then-import; hard to reason about.
- Hidden side effects make testing painful (can't disable, can't substitute).
- Init failures must `panic` — no graceful handling.

Prefer explicit constructors called from `main`. Use `init` only for things
that genuinely can't fail: registering encoders, computing constants.

---

### Q12. What's the deal with `iota`?

Compile-time auto-incrementing identifier inside `const` blocks:

```go
const (
    Red   = iota // 0
    Green        // 1
    Blue         // 2
)
```

Reset to 0 in each const block. Useful for enums, bit flags
(`1 << iota`), and skip values (`_ = iota`).

Trap: `iota` increments **per line in the const block**, not per use.

---

### Q13. What does `_ = imported.Package` do?

Imports the package **for its side effects** (init() runs) without using
any exported names. Common for driver registration:

```go
import _ "github.com/lib/pq"  // registers postgres driver with database/sql
```

---

### Q14. Two strings, same content, allocated separately — are they `==`?

Yes. String equality is **by value** (compares contents). Interfaces holding
those strings also compare equal. Slices/maps/funcs are NOT comparable —
trying to put one in a `map[any]X` panics at runtime if the key is non-comparable.

---

### Q15. What's the zero value of a struct with mixed fields?

Each field gets its own zero value: numbers → 0, strings → "", bools → false,
pointers/maps/slices/channels/funcs/interfaces → nil. A zero-valued struct
is always valid memory — many Go types are designed to be useful at zero
value (`sync.Mutex`, `bytes.Buffer`, `strings.Builder`).
