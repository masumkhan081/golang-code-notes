# Testing

---

### Q1. What does an idiomatic Go test look like?

Table-driven, with subtests:

```go
func TestAdd(t *testing.T) {
    cases := []struct {
        name    string
        a, b, want int
    }{
        {"positives", 1, 2, 3},
        {"with zero", 0, 5, 5},
        {"negatives", -1, -1, -2},
    }
    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            if got := Add(tc.a, tc.b); got != tc.want {
                t.Errorf("Add(%d,%d)=%d, want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}
```

---

### Q2. `t.Error` vs `t.Fatal`?

- `t.Error` — record failure, **continue** the test.
- `t.Fatal` — record failure, **stop** the test (calls `runtime.Goexit`).

Use Fatal when later assertions depend on the failed step (e.g., setup failed, no point checking results).

Don't call `t.Fatal` from a goroutine other than the test's — it won't work.
Send a result back via channel and assert in the main goroutine.

---

### Q3. How does `t.Parallel()` work?

Marks the test (or subtest) as eligible to run in parallel with **other parallel tests**.
- Tests without `t.Parallel()` run sequentially first.
- Then all parallel tests run concurrently, capped by `-parallel=N` (default GOMAXPROCS).

Trap: in a `for _, tc := range cases` loop **before Go 1.22**, you must
shadow `tc := tc` before `t.Run` or all subtests see the last value.

---

### Q4. What's `TestMain` for?

Package-level setup/teardown. If defined, the testing harness calls it
instead of running tests directly:

```go
func TestMain(m *testing.M) {
    db = openTestDB()
    code := m.Run()
    db.Close()
    os.Exit(code)
}
```

Use for: DB connection, mock server, fixture seeding, leak detection
(`goleak.VerifyTestMain`).

Gotcha: `defer` doesn't run after `os.Exit`. Do teardown before, or use
a `func() int { ... }` wrapper.

Code: `10_testing/testmain_test.go`.

---

### Q5. How do you write a benchmark and read its output?

```go
func BenchmarkConcat(b *testing.B) {
    parts := setup()
    b.ResetTimer()    // exclude setup
    b.ReportAllocs()  // print B/op and allocs/op
    for i := 0; i < b.N; i++ {
        _ = concat(parts)
    }
}
```

Run: `go test -bench=. -benchmem -count=10`

Output:
```
BenchmarkConcat-8   100000   12345 ns/op   2048 B/op   3 allocs/op
                ^^^ ^^^^^^   ^^^^^^^^^^^   ^^^^^^^^^   ^^^^^^^^^^^^
              CPUs    N      per-op time   bytes/op    allocs/op
```

For comparing two versions, use `benchstat`:
```
go install golang.org/x/perf/cmd/benchstat@latest
benchstat old.txt new.txt
```

Code: `10_testing/benchmark_advanced_test.go`.

---

### Q6. What is fuzzing and what should you fuzz?

Coverage-guided property-based testing (Go 1.18+).

```go
func FuzzReverse(f *testing.F) {
    f.Add("hello")  // seed
    f.Fuzz(func(t *testing.T, s string) {
        if Reverse(Reverse(s)) != s {
            t.Errorf("round-trip failed for %q", s)
        }
    })
}
```

Run: `go test -run=^$ -fuzz=FuzzReverse -fuzztime=10s`

Fuzz **parsers, decoders, anything that takes []byte/string and shouldn't
panic, round-trip pairs (encode/decode), and invariants**. Failing inputs
are saved to `testdata/fuzz/` and become permanent regression tests.

Code: `10_testing/fuzz_test.go`.

---

### Q7. How do you mock dependencies in Go?

- **Interfaces at the consumer** + a hand-written fake for tests. Idiomatic.
- `httptest.NewServer` for HTTP clients.
- `testing/fstest.MapFS` for filesystem.
- `database/sql/driver` fakes (or libraries like `sqlmock`) — but prefer real DB in containers when possible.
- Generated mocks (`mockery`, `gomock`) when the interface is large.

Strong opinion in the community: keep mocks small, prefer fakes over mocks,
prefer integration tests with real dependencies in Docker for the hot
paths.

---

### Q8. How do you test HTTP handlers?

`net/http/httptest`:

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
    rr := httptest.NewRecorder()
    Handler(rr, req)
    if rr.Code != http.StatusOK {
        t.Errorf("status=%d, want 200", rr.Code)
    }
}
```

For testing a **client**, spin up a server:
```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"ok":true}`))
}))
defer ts.Close()
// pass ts.URL to your client
```

---

### Q9. What flags do you always run in CI?

```bash
go vet ./...
go test -race -covermode=atomic -coverprofile=cover.out ./...
```

- `-race` finds data races
- `-covermode=atomic` is required when combining `-race` and coverage
- `vet` catches lost cancels, copied locks, shadowed errors

Detail: `10_testing/race_and_coverage.md`.

---

### Q10. What's `t.Helper()` and `t.Cleanup()`?

**`t.Helper()`** — call at the top of an assertion helper so that when it
fails, the failure line points to the **caller**, not inside the helper:

```go
func assertEqual(t *testing.T, got, want int) {
    t.Helper() // failure shows the caller's line number
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}
```

Without it, every test failure points to the same line inside `assertEqual`
— useless.

**`t.Cleanup(fn)`** — registers a teardown that runs when the test (and all
its subtests) finishes. Runs in LIFO order, even on `t.Fatal`.

```go
func TestThing(t *testing.T) {
    f, err := os.CreateTemp("", "x")
    if err != nil { t.Fatal(err) }
    t.Cleanup(func() { os.Remove(f.Name()) }) // safer than defer
}
```

Beats `defer` because:
- Runs after subtests too (defer runs when the outer test function returns,
  but subtests can outlive that scope in some setups).
- Survives `t.Fatal` — defer also does, but Cleanup composes nicely when
  helpers register their own teardowns.

A helper that creates a temp resource can register its own cleanup,
hiding setup AND teardown from the test:

```go
func newTestDB(t *testing.T) *DB {
    t.Helper()
    db := openTestDB()
    t.Cleanup(func() { db.Close() })
    return db
}
```

---

### Q11. What's a golden file test?

Compare actual output to a checked-in expected file (the "golden" file).
Useful for snapshot-like assertions (templates, generated code, formatted output).

```go
got := render(input)
golden := filepath.Join("testdata", t.Name()+".golden")
if *update {
    os.WriteFile(golden, got, 0644)
}
want, _ := os.ReadFile(golden)
if !bytes.Equal(got, want) {
    t.Errorf("mismatch with %s", golden)
}
```

The `-update` flag (a custom `flag.Bool`) regenerates goldens after
intentional changes — review the diff before committing.
