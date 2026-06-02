# 10_testing

Unit tests, table-driven tests, benchmarks, fuzzing, race & coverage.

## Contents

- `main_test.go` — basic table-driven test + benchmark
- `subtests_test.go` — `t.Run` subtests
- `http_handler_test.go` — `httptest` for HTTP handlers
- `testmain_test.go` — **`TestMain` package-level setup/teardown**
- `benchmark_advanced_test.go` — **`ResetTimer`, `ReportAllocs`, sub-benchmarks, `RunParallel`, benchstat**
- `fuzz_test.go` — **native fuzzing (Go 1.18+) with seed corpus**
- `race_and_coverage.md` — **`-race` and `-cover` CLI reference**

## Cross-references

- Interview Q&A: `../25_interview_questions/04_testing.md`
- Concurrency leak detection (use in `TestMain` with goleak): `../11_concurrency/08_goroutine_leaks/`
