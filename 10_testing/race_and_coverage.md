# Race detector and coverage

CLI knowledge interviewers expect you to recite without thinking.

## Race detector

```bash
go test -race ./...        # always run tests with -race in CI
go run -race main.go       # spot races during development
go build -race -o app .    # ship a race-instrumented build only for staging
```

What it catches: **data races** — two goroutines accessing the same memory,
at least one of them writing, without a happens-before edge (mutex, channel,
sync primitive, or atomic).

What it does not catch:
- Races on memory that never happen to be touched concurrently during the run.
- Logic races (e.g., TOCTOU) where each access is properly synchronized
  but the sequence is wrong.
- Race-free code with deadlocks.

Cost: roughly **2-20x slower**, **5-10x more memory**. Don't ship `-race`
binaries to production except temporary diagnosis.

Sample output:
```
WARNING: DATA RACE
Read at 0x00c0000160e8 by goroutine 7:
  main.racy.func1()
      main.go:14 +0x44
Previous write at 0x00c0000160e8 by goroutine 6:
  main.racy.func1()
      main.go:14 +0x58
```

The fix is always: synchronize the access (Mutex, channel, atomic), or
prove non-overlap (one goroutine writes, then sends via a channel, then
the other reads).

## Coverage

```bash
go test -cover ./...                          # summary % per package
go test -coverprofile=cover.out ./...         # write profile
go tool cover -func=cover.out                 # per-function coverage
go tool cover -html=cover.out                 # browser heatmap
go tool cover -html=cover.out -o cover.html   # save heatmap
```

By default coverage is per-package — only code in the same package as the
test file counts. To cover code from other packages (e.g., integration
tests in `_test` packages):

```bash
go test -coverpkg=./... ./...
```

Coverage modes (`-covermode=`):
- `set` — was each line hit? (default, cheapest)
- `count` — how many times was it hit?
- `atomic` — like count but safe under `-race` (use this with concurrent tests)

```bash
go test -race -covermode=atomic -coverprofile=cover.out ./...
```

## Putting it together in CI

```bash
go vet ./...
go test -race -covermode=atomic -coverprofile=cover.out ./...
go tool cover -func=cover.out | tail -1   # total coverage line
```

## Useful flags reference

| Flag | Purpose |
|---|---|
| `-race` | enable race detector |
| `-cover` / `-coverprofile=f` | enable coverage / dump profile |
| `-covermode=atomic` | coverage counter mode (use with -race) |
| `-coverpkg=./...` | include code from other packages |
| `-short` | tests should call `testing.Short()` and skip slow ones |
| `-run=Regex` | run only matching tests (`-run=^$` disables all) |
| `-bench=Regex` | run benchmarks (`-bench=.` for all) |
| `-benchtime=5s` / `1000x` | benchmark duration or iterations |
| `-count=N` | repeat each test/benchmark N times |
| `-cpu=1,2,4` | sweep GOMAXPROCS |
| `-v` | verbose, prints t.Log output |
| `-failfast` | stop on first test failure |
| `-timeout=30s` | kill the binary after duration (default 10m) |
| `-parallel=N` | max parallel `t.Parallel()` tests |
