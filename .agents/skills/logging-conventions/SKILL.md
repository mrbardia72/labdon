---
name: logging-conventions
description: Apply this repo's structured logging conventions — log levels, required fields (request/trace id, layer, operation), and what belongs in a log line versus a returned error. Use when adding or editing any log.* call in Go code.
---

# Logging Conventions

This skill defines what a log line should contain and when to emit one. It
pairs with `error-handling-conventions`'s "log vs return" rule — that skill
says *when* to log; this one says *what* the log call should look like once
you've decided to.

## Structured fields, not formatted strings

Every log call uses key-value fields, never a pre-formatted string that
embeds variable data via `fmt.Sprintf`. The message itself stays a short,
static, lowercase phrase.

```go
// ❌ Avoid — unparseable, inconsistent field names across call sites
log.Info(fmt.Sprintf("created widget %s for owner %s", widget.ID, ownerID))

// ✅ Prefer
log.Info(ctx, "widget created",
	"widget_id", widget.ID.String(),
	"owner_id", ownerID.String(),
)
```

- Field names are `snake_case` and stable across the codebase — `widget_id`
  everywhere a widget ID is logged, not `id` in one place and `widgetID` in
  another. If a field for this concept already exists elsewhere in the repo,
  reuse its exact name.
- Always pass `ctx` as the first argument when the logger supports it — this
  is how request/trace correlation fields get attached automatically. Don't
  hand-add a `request_id` field if the logger already injects one from
  context.

## Levels

| Level | Use for |
|---|---|
| `Debug` | Internal state useful only when actively diagnosing a problem; off by default in production. |
| `Info` | Notable things that happened successfully (resource created, job completed, server started) — not routine per-request noise. |
| `Warn` | Something recoverable happened that a human should eventually notice (fallback used, retry succeeded after N attempts, deprecated path hit). |
| `Error` | A terminal failure at the point it's logged — see `error-handling-conventions` for when a failure is "terminal" versus something that should just be returned instead. |

Don't log the same failure at `Error` in three layers as it propagates — log it
once, at the layer where it's terminal (becomes a response, ends a job, or is
otherwise not going to be wrapped further).

## Required fields by context

- **HTTP request lifecycle** (via middleware, not per-handler): method, path,
  status, duration, request id. Individual handlers don't re-log this.
- **Business operation logs** (when an extension like `business-layer-extensions`
  adds a logging concern): the operation name, the primary entity ID(s)
  involved, and outcome (`status: success|error`) — matching the metrics
  extension's `status` label convention so logs and metrics stay correlated.
- **Background jobs/cron**: job name, start, completion or failure, and item
  counts processed — enough to answer "did this run, and did it do anything?"
  from logs alone without needing to attach a debugger.

## Error fields

When logging an error (per the `Error` level above), pass the error as its own
field, not interpolated into the message:

```go
log.Error(ctx, "widget creation failed",
	"error", err,
	"owner_id", ownerID.String(),
)
```

Log the fully wrapped error (after `%w` wrapping has added context per
`error-handling-conventions`), so the log line carries the whole causal chain,
not just the innermost message.

## What NOT to log

- Secrets, tokens, full request bodies containing PII, or raw passwords —
  even at `Debug`. If a struct being logged might contain one of these, log
  specific safe fields instead of the whole struct.
- The same event twice because both a logging extension and the handler log
  it independently — see `business-layer-extensions` for how the
  decorator pattern is supposed to centralize this concern in one place
  instead of scattering log calls through core business logic.
- Routine successful reads (`Get`, `List`) at `Info` in a hot path — that
  belongs at `Debug` if it's needed at all, otherwise it drowns out the
  `Info` logs that actually matter.

## Extension-layer logging

When logging is added as a cross-cutting concern via the `ExtBusiness`/
`Extension` seam (see `business-layer-extensions`), the logging extension wraps
each method, logs outcome using the field conventions above, and returns the
original result unchanged — it never swallows or alters the error it's
logging.