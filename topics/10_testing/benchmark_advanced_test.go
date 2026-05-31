// benchmark_advanced_test.go
// Going beyond `for i := 0; i < b.N; i++`.
//
// Run:
//   go test -bench=.                    // all benchmarks
//   go test -bench=BenchmarkConcat      // one benchmark
//   go test -bench=. -benchmem          // also report allocs
//   go test -bench=. -benchtime=5s      // run each for 5s
//   go test -bench=. -benchtime=1000x   // run each exactly 1000 times
//   go test -bench=. -count=10          // repeat 10x for variance
//   go test -bench=. -cpu=1,2,4,8       // sweep GOMAXPROCS
package main

import (
	"strings"
	"testing"
)

// Two ways to concatenate strings — common benchmark candidate.
func concatPlus(parts []string) string {
	var s string
	for _, p := range parts {
		s += p // O(n²) — allocates a new string each step
	}
	return s
}

func concatBuilder(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

var input = strings.Split(strings.Repeat("abc ", 100), " ")

// BAD: setup is counted in the benchmark time.
func BenchmarkBadSetup(b *testing.B) {
	parts := strings.Split(strings.Repeat("abc ", 1000), " ") // measured!
	for i := 0; i < b.N; i++ {
		_ = concatBuilder(parts)
	}
}

// GOOD: ResetTimer after setup so only the hot path is measured.
func BenchmarkConcatBuilder(b *testing.B) {
	parts := strings.Split(strings.Repeat("abc ", 1000), " ")
	b.ResetTimer() // exclude setup
	b.ReportAllocs() // print B/op and allocs/op
	for i := 0; i < b.N; i++ {
		_ = concatBuilder(parts)
	}
}

func BenchmarkConcatPlus(b *testing.B) {
	parts := strings.Split(strings.Repeat("abc ", 1000), " ")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = concatPlus(parts)
	}
}

// Sub-benchmarks — like sub-tests, scoped names + isolated metrics.
func BenchmarkConcatSizes(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		parts := strings.Split(strings.Repeat("abc ", size), " ")
		b.Run("builder/"+itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = concatBuilder(parts)
			}
		})
	}
}

// Parallel benchmarks — RunParallel splits b.N across PB goroutines.
// Use for measuring contention (mutex, channel, atomic).
func BenchmarkParallelBuilder(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = concatBuilder(input)
		}
	})
}

// Pause the timer during expensive per-iteration setup.
func BenchmarkWithPauseResume(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		parts := strings.Split(strings.Repeat("x ", 100), " ") // not measured
		b.StartTimer()
		_ = concatBuilder(parts) // measured
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

/*
READING THE OUTPUT
  BenchmarkConcatBuilder-8   100000   12345 ns/op   2048 B/op   3 allocs/op
                       ^^^   ^^^^^^   ^^^^^^^^^^^   ^^^^^^^^^   ^^^^^^^^^^^^
                       GOMAXPROCS  N  per-op time   bytes/op    allocs/op

COMPARING TWO COMMITS — benchstat
  go install golang.org/x/perf/cmd/benchstat@latest
  go test -bench=. -count=10 > old.txt
  # ... change code ...
  go test -bench=. -count=10 > new.txt
  benchstat old.txt new.txt

PROFILING DURING BENCHMARKS
  go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out
  go tool pprof cpu.out

DON'TS
  - Don't `time.Sleep` in a benchmark — it inflates ns/op meaninglessly.
  - Don't forget to USE the result (assign to a package-level sink) — the
    compiler will dead-code-eliminate work that produces unused values.
  - Don't trust a single run. Use -count=10 + benchstat for any claim
    smaller than ~10%.
*/
