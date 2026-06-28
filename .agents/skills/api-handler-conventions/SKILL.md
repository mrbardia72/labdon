---
name: api-handler-conventions
description: Apply this repo's App-layer HTTP handler conventions — request decode/validate/toBus flow, status code mapping, middleware chain shape, and response envelope. Use when writing or editing handler code under app/domain/*/*.go or App-layer routing/middleware files.
---

# API Handler Conventions

This skill defines the standard shape of an HTTP handler in the App layer
(`app/domain/<x>app`). It complements `layered-architecture-types` (which governs
the request/response *types*) by covering the handler *flow* around those types:
decode, validate, call Business, map result, respond.

## Handler shape

Every handler follows the same four-step shape, in this order:

```go
func (h *Handlers) Create(ctx context.Context, r *http.Request) web.Encoder {
	var req CreateWidgetRequest
	if err := web.Decode(r, &req); err != nil {
		return errs.NewFieldsError("decode", err)
	}

	input, err := toBusCreateWidget(req)
	if err != nil {
		return err // already errs.FieldErrors or a typed validation error
	}

	widget, err := h.widgetBus.Create(ctx, input)
	if err != nil {
		return errs.New(errs.Internal, fmt.Errorf("create: %w", err))
	}

	return fromBusCreateWidgetResponse(widget)
}
```

1. **Decode** — only `web.Decode`/`json.Unmarshal` into the primitive request
   struct. No business logic here.
2. **Validate + convert** — call the named `toBus<Type>` converter (see
   `layered-architecture-types`). Validation errors return directly; don't wrap
   `errs.FieldErrors` in another error type.
3. **Call Business** — exactly one Business-layer call per handler when
   possible. If a handler must call two Business domains, that's expected at the
   App layer (App may import multiple Business packages — see
   `layered-architecture-types`'s import rules) — just keep each call's error
   handling distinct so you know which one failed.
4. **Convert + respond** — `fromBus<Type>Response`, never return a Business
   model or DB row type directly from a handler.

## Status code mapping

Map Business sentinel errors to status codes with `errors.Is`, in one place per
handler (or one shared helper if the same mapping repeats across handlers in a
domain) — never infer status from a message string.

```go
switch {
case errors.Is(err, widgetbus.ErrNotFound):
	return errs.New(errs.NotFound, err)
case errors.Is(err, widgetbus.ErrAlreadyExists):
	return errs.New(errs.AlreadyExists, err)
case errors.As(err, &errs.FieldErrors{}):
	return err // FieldErrors already encodes its own status
default:
	return errs.New(errs.Internal, err)
}
```

See `error-handling-conventions` for how sentinels are declared and wrapped
before they reach this point.

## Middleware chain

- Middleware order (outermost to innermost) is: panic-recovery → request-id/
  logging → auth → handler. Don't reorder without a concrete reason — recovery
  must wrap everything so a panic anywhere still produces a clean 5xx.
- A middleware that short-circuits the chain (auth rejection, rate limit) writes
  the response itself and returns — it does not call `next` and then also let
  the handler run.
- New middleware is added as its own small function (`func(web.Handler)
  web.Handler` or this repo's equivalent signature), composed in the route
  registration — not inlined into an existing middleware's body.

## Route registration

- One `routes.go` (or equivalent) per App domain package registers that
  domain's endpoints; it doesn't reach into another domain's routes (consistent
  with the no-cross-domain-App-import rule in `layered-architecture-types`).
- Path naming: plural resource nouns, lowercase, hyphenated for multi-word
  resources (`/widget-groups`, not `/widgetGroups` or `/widget_groups`).
- Use Go 1.22+ enhanced `http.ServeMux` patterns (`mux.HandleFunc("POST
  /widgets/{id}", handler)`, `r.PathValue("id")`) per `use-modern-go` — don't
  hand-roll path parsing.

## Pagination, filtering, sorting

If a list endpoint supports paging, the request struct carries primitive
`page`/`rows` (or this repo's established names) with documented defaults
applied in `toBus<Type>` (default if zero/absent, validate if present and
invalid) rather than in the handler body.

## What NOT to do

- Don't put `business/types` strong types in the request/response struct
  (covered fully by `layered-architecture-types`, but it's the single most
  common handler mistake, so it's worth restating here).
- Don't call the Storage layer directly from a handler — always through
  Business.
- Don't swallow the decode error to produce a generic 400 — return the
  validation detail so the API caller can fix the request (see
  `error-handling-conventions`'s `errs.FieldErrors` section).
- Don't log inside the handler for errors that are already being returned and
  will be logged once at the middleware/response-writing boundary — see the
  "log vs return" rule in `error-handling-conventions`.