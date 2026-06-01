// fuzz_test.go
// Native fuzzing (Go 1.18+). Run with:
//   go test -run=^$ -fuzz=FuzzReverse -fuzztime=10s
//
// -run=^$  disables unit tests so only the fuzz target runs.
// -fuzz    selects the fuzz target by name.
// -fuzztime is how long to keep generating inputs (default: until you Ctrl-C).
package main

import (
	"testing"
	"unicode/utf8"
)

// Reverse a UTF-8 string by runes (not bytes — that's the bug magnet).
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func FuzzReverse(f *testing.F) {
	// Seed corpus — fuzz engine mutates these.
	// Also acts as a regression suite: any failing input found by the
	// engine is auto-saved to testdata/fuzz/<TestName>/ and replayed
	// as a normal test on subsequent runs.
	for _, seed := range []string{"hello", "Go", "", "a", "日本語", "🙂"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		// Property #1: reversing twice yields the original.
		got := Reverse(Reverse(in))
		if got != in {
			t.Errorf("double-reverse mismatch: %q -> %q", in, got)
		}
		// Property #2: output stays valid UTF-8 if input was.
		if utf8.ValidString(in) && !utf8.ValidString(Reverse(in)) {
			t.Errorf("Reverse produced invalid UTF-8 from %q", in)
		}
	})
}

/*
FUZZING WORKFLOW
  1. Write a property — something that must hold for ALL valid inputs
     (round-trip, idempotence, no panic, invariant preserved).
  2. Seed with f.Add(...) — known interesting cases.
  3. Run with -fuzz. The engine mutates seeds via coverage-guided search.
  4. On a failure, Go writes the minimized input to
     testdata/fuzz/FuzzXxx/<hash>. Commit this file — it becomes a
     permanent regression test that runs under `go test` (no -fuzz needed).

WHAT TO FUZZ
  - Parsers, decoders, deserializers
  - Anything that takes []byte or string and must not panic
  - Round-trip pairs: encode/decode, marshal/unmarshal, compress/decompress
  - State machines: sequence of operations should not deadlock

WHAT NOT TO FUZZ
  - Pure arithmetic that's already exhaustively unit-tested
  - I/O-bound code (fuzz the parsing layer instead)

SUPPORTED INPUT TYPES
  string, []byte, bool, byte, rune, float32/64,
  int/int8/16/32/64, uint/uint8/16/32/64
  Multiple parameters allowed: f.Fuzz(func(t *testing.T, a, b int){...})
*/
