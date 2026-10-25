# Copilot instructions

Repo-wide guidance for GitHub Copilot. The full agent contract is in
[AGENTS.md](AGENTS.md); the per-topic conventions are in `.agents/skills/`.

## Project
- Language: **Go** (toolchain pinned in `go.mod`). No TypeScript/JS code.
- Layout: `cmd/` entrypoints, `app/` HTTP layer, `business/` domain logic,
  `foundation/` shared infrastructure (logger, trace, database, buildinfo).
- Architecture & type-boundary rules: see ARCHITECTURE.md and the
  `layered-architecture-types` skill. Primitives at the edges, strong types
  in `business/types/*`, named converters (`toBus<Type>`, `fromBus<Type>Response`,
  `toDB<Type>`) at every layer crossing.

## When writing code
- Follow `error-handling-conventions` (wrap with `%w`, sentinels, `errs.FieldErrors`).
- Follow `logging-conventions` (structured zap fields; never log secrets).
- Follow `go-testing-conventions` (table-driven tests, `t.Context()`, error-path assertions).
- Prefer modern Go features enabled by the `go.mod` version (`use-modern-go`).

## When reviewing / suggesting
- Do not invent APIs for packages that don't exist yet; scaffold minimal code.
- Never suggest committing, pushing, or git mutations unprompted (AGENTS.md rule).

## Commands
- `make help` — all targets. `make check` — everything CI runs.
- `make test-only` — fast local tests; `make lint` requires `.golangci.yml` (present).
