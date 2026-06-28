---
name: error-handling-conventions
description: Apply this repo's Go error-handling conventions — wrapping with %w, sentinel/typed errors per domain, errs.FieldErrors for validation, and where errors get logged vs returned. Use when writing or editing any Go code that returns, wraps, or handles an error, not only when explicitly asked.
---

# Error Handling Conventions

This skill defines **how to produce and propagate** errors in this repo's Go
code. It is the authoring counterpart to the `error-tripwire` review lens, which
only judges whether a finished diff's failure paths are visible. Apply this
whenever writing or editing code that returns, wraps, logs, or handles an `error`.

## Wrap with context, always

Every error returned up a call boundary gets wrapped with `fmt.Errorf("...: %w",
err)` describing the operation that failed — never the bare `return err` once
you've crossed a meaningful boundary (e.g. Storage → Business, Business → App).

```go
// ❌ Avoid — caller has no idea which query failed
func (s *Store) GetByID(ctx context.Context, id string) (dbWidget, error) {
	var row dbWidget
	if err := s.db.GetContext(ctx, &row, query, id); err != nil {
		return dbWidget{}, err
	}
	return row, nil
}

// ✅ Prefer
func (s *Store) GetByID(ctx context.Context, id string) (dbWidget, error) {
	var row dbWidget
	if err := s.db.GetContext(ctx, &row, query, id); err != nil {
		return dbWidget{}, fmt.Errorf("get widget %s: %w", id, err)
	}
	return row, nil
}
```

Within a single function body, intermediate `return err` between two lines of the
same operation is fine — wrap once per boundary crossing, not once per call.

## Sentinel and typed errors live in the Business layer

Domain-meaningful errors (`ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidState`,
...) are declared as package-level `var Err... = errors.New("...")` in the
relevant `business/domain/<x>bus` package, next to the model type. Storage and App
layers reference these, they don't invent their own.

- Storage maps driver-specific failures (`sql.ErrNoRows`, unique-violation codes)
  into the Business sentinel at the Storage→Business boundary — never let
  `sql.ErrNoRows` leak past the store.
- App maps Business sentinels to HTTP status via `errors.Is`, not by inspecting
  error message text.

```go
// store: translate driver error to the domain sentinel
if errors.Is(err, sql.ErrNoRows) {
	return widgetbus.Widget{}, fmt.Errorf("get widget %s: %w", id, widgetbus.ErrNotFound)
}

// app handler: map sentinel to status
switch {
case errors.Is(err, widgetbus.ErrNotFound):
	return response.NotFound(...)
case errors.Is(err, widgetbus.ErrAlreadyExists):
	return response.Conflict(...)
default:
	return response.InternalError(...)
}
```

For errors that carry data (not just identity), define a typed error
(`type ValidationError struct { Field string; Reason string }`) implementing
`error`, and match with `errors.As` (or `errors.AsType[T]` on Go 1.26+, per
`use-modern-go`).

## Validation errors use errs.FieldErrors

Field-level validation in App-layer `toBus<Type>` converters (see
`layered-architecture-types`) accumulates into `errs.FieldErrors`, not a single
generic error — this lets the API response surface which field failed.

```go
var fieldErrors errs.FieldErrors

ownerID, err := types.ParseOwnerID(req.OwnerID)
if err != nil {
	fieldErrors.Add("owner_id", err)
}
name := strings.TrimSpace(req.Name)
if name == "" {
	fieldErrors.Add("name", errors.New("required"))
}

if len(fieldErrors) > 0 {
	return widgetbus.CreateWidgetInput{}, fieldErrors
}
```

Never short-circuit on the first field error in a multi-field request — collect
all of them, so the caller fixes the request in one round trip.

## Logging vs returning — pick one per error

An error is either **returned** (caller decides what to do) or **logged at the
point it's terminal** (nothing upstream can act on it) — never both, and never
neither.

- **Don't log-and-return.** Logging an error and also returning it usually means
  it gets logged twice (once here, once by whoever logs at the boundary) or the
  log has no actionable context the wrapped error doesn't already carry.
- **Do log at genuinely terminal points**: a goroutine with no channel to report
  back, a `defer`'d cleanup whose failure can't propagate, a background job's
  final retry exhaustion.
- Background goroutines spawned with `wg.Go` (Go 1.25+, per `use-modern-go`)
  still need an explicit way to surface failure — a buffered error channel, an
  `errgroup`, or a logged terminal error. A goroutine that silently drops its
  error is the single most common G4 finding from `error-tripwire`.

## Defer and cleanup errors

If a deferred `Close()`/`Rollback()` can fail meaningfully, don't discard it with
a bare `defer rows.Close()` when that failure matters (e.g. a write-path
transaction). Capture it:

```go
defer func() {
	if cerr := tx.Rollback(); cerr != nil && !errors.Is(cerr, sql.ErrTxDone) {
		log.Error(ctx, "rollback failed", "error", cerr)
	}
}()
```

For read-only resources where a close failure is not actionable (e.g.
`defer rows.Close()` after a successful scan loop), the bare form is fine —
don't manufacture handling for a failure mode with no consequence.

## Context cancellation

Check `ctx.Err()` (or let the underlying call do it) before starting expensive
work in a loop; don't swallow `context.Canceled` / `context.DeadlineExceeded` as
if they were ordinary errors — propagate them with `%w` like any other error so
callers can `errors.Is(err, context.Canceled)` if they need to distinguish
"caller gave up" from "operation failed."

## What NOT to do

- Don't introduce a new sentinel error for something `errors.Is(err,
  sql.ErrNoRows)` or an existing Business sentinel already covers.
- Don't use `panic`/`recover` for ordinary failure handling — reserve `panic` for
  programmer errors (invariant violations that indicate a bug), and only
  `recover` at a process-boundary (e.g. HTTP middleware) to convert it into a
  5xx response plus a logged stack trace.
- Don't return `nil` error alongside a zero-value/partial result to mean
  "not found" — return the sentinel error instead, so callers can't
  accidentally treat a missing value as a valid empty one.