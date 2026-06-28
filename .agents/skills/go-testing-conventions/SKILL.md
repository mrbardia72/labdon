---
name: go-testing-conventions
description: Apply this repo's Go test-writing conventions — table-driven tests, subtest naming, t.Context() usage, error-path assertions with errors.Is/As, and layer-crossing test shape. Use when writing or editing any _test.go file, not only when explicitly asked for test guidelines.
---

# Go Testing Conventions

This skill defines **how to write** a Go test in this repo. It is the authoring
counterpart to the `harness-map` review lens (which only judges whether existing
tests give good regression coverage). Apply this whenever creating or editing a
`_test.go` file.

## Table-driven tests are the default

Any function with more than one meaningful input/output pair gets a table, not a
sequence of separate `Test*` functions.

```go
func TestParseWidgetID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    types.WidgetID
		wantErr bool
	}{
		{name: "valid uuid", input: validID.String(), want: validID},
		{name: "empty string", input: "", wantErr: true},
		{name: "malformed uuid", input: "not-a-uuid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types.ParseWidgetID(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseWidgetID(%q): want error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseWidgetID(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseWidgetID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
```

- Subtest names are short, lowercase, hyphen/space-separated descriptions of the
  *case*, not the assertion (`"empty string"`, not `"returns error"`).
  - Every case needs a `name` field — never rely on table index for identification.
- Group success and failure cases in the same table; don't split into a
  `TestXSuccess` / `TestXFailure` pair unless setup genuinely diverges.

## Context in tests

ALWAYS use `t.Context()` (Go 1.24+) instead of manually building and cancelling a
context. See `use-modern-go` for the general rule; this is the test-specific
reminder since it comes up in nearly every Business/Storage test.

```go
func TestBus_Create(t *testing.T) {
	ctx := t.Context()
	// ...
}
```

Only fall back to `context.WithTimeout`/`WithCancel` when the test needs to
exercise cancellation/timeout behavior itself — that's the one case `t.Context()`
doesn't cover.

## Error-path assertions

Never compare errors with `==` or `.Error()` string matching when the underlying
error is a sentinel or typed error. Use `errors.Is` / `errors.As` (or
`errors.AsType[T]` once the module targets Go 1.26+, per `use-modern-go`).

```go
// ❌ Avoid — breaks the moment the message wording changes
if err.Error() != "widget not found" { ... }

// ✅ Prefer
if !errors.Is(err, widgetbus.ErrNotFound) {
	t.Fatalf("want ErrNotFound, got %v", err)
}
```

`.Error()` string matching is acceptable only for a one-off assertion where no
sentinel/typed error exists and the message itself is the contract being tested
(rare — confirm there's truly no typed error first).

## What a test must cover per changed behavior

For each new or modified exported function/method, at minimum:

1. **Happy path** — the common valid input.
2. **Validation/boundary failure** — the input that should be rejected, asserted
   via `errors.Is`/`errors.As` or `errs.FieldErrors` inspection, not just "an
   error was returned."
3. **Layer-crossing correctness**, if the function is a converter
   (`toBus<Type>`, `fromBus<Type>Response`, `toDB<Type>`) — assert the specific
   field mapping, not just that the call didn't panic.

This is the same bar `harness-map` will check for in review — writing to this
standard up front avoids G3/G4 findings later.

## Test doubles and mocking

- Prefer a small hand-written fake over a generated mock when the interface is
  narrow (consumer-owned interfaces are usually 1-3 methods in this repo — see
  `layered-architecture-types`).
- A fake's behavior must be controllable per-test-case (inject the error/value
  through the table), not hardcoded — otherwise it silently drifts from being a
  regression guard into a tautology.
- Don't assert against a mock's internal call count unless call count *is* the
  behavior under test (e.g. verifying retry-exhaustion stops at N attempts).

## Race and parallelism

- Tests in this repo run with `-race` in CI (see `go-basic-checks` /
  `test` jobs) — any new goroutine spawned by code under test, or by the test
  itself, must be safe under `-race` and must complete (or be explicitly
  cancelled via `t.Context()`) before the test returns.
- Use `t.Parallel()` inside each `t.Run` subtest when subtests are independent
  and don't share mutable fixtures — but never mark a subtest parallel if its
  fake/fixture is shared and mutated across cases.

## Naming and location

- Test file: `<file>_test.go` next to the file it tests, package `<pkg>` (white-box)
  unless the test is asserting only the public API, in which case `<pkg>_test`
  (black-box) is preferred to catch accidental reliance on internals.
- Test function name: `Test<Type>_<Method>` for methods, `Test<FuncName>` for
  free functions.

## Benchmarks

Use `b.Loop()` (Go 1.24+), not `for i := 0; i < b.N; i++` — see `use-modern-go`.
Benchmarks live alongside the test file (`Benchmark<Name>` in the same
`_test.go`) unless the package has enough of them to warrant a dedicated
`<file>_bench_test.go`.

## When NOT to apply

- Generated code and its tests (skip; fix the generator).
- Pure integration/e2e suites with their own established conventions — this
  skill covers unit/package-level Go tests under `app/*`, `business/*`,
  `foundation/*`.