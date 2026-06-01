# 25_interview_questions

Topic-grouped Go interview Q&A. Each file is standalone; answers point to
runnable code elsewhere in the repo.

## Contents

- `01_concurrency.md` — goroutines, channels, sync primitives, GMP, race vs deadlock
- `02_interfaces.md` — satisfaction, typed-nil, method sets, accept-iface-return-struct
- `03_context.md` — cancellation, timeouts, WithValue, propagation
- `04_testing.md` — table-driven, subtests, parallel, benchmarks, fuzzing, mocking
- `05_misc_traps.md` — slice aliasing, map iteration, defer evaluation, errors.Is, nil receivers
- `06_generics.md` — constraints, `~T`, type inference, method-parameter limitation, stdlib generic packages

## How to use

- Read the answer, then open the linked code and run it.
- Try to answer aloud before reading — gaps surface fast that way.
- For each "gotcha" question, write your own minimal repro from memory.
