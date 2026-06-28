# `.agents/skills`

Skills for AI coding agents (Amp, Claude Code, etc.) working in this repo.
Each `SKILL.md` has a `description` the agent uses to decide relevance — no
manual invocation needed in most cases.

This directory is meant to be copied wholesale into any new Go project that
follows this repo's conventions. Two categories:

## Authoring skills (apply while writing code)

| Skill | Triggers on |
|---|---|
| `use-modern-go` | Any `.go` file — syntax/feature guidance gated on the project's `go.mod` version. |
| `layered-architecture-types` | `app/*`, `business/domain/*`, `.../stores/*db` — primitive/strong-type boundary + converters. |
| `business-layer-extensions` | Adding tracing/logging/metrics/caching/auth to a business domain via the `ExtBusiness`/`Extension` seam. |
| `branching-logic-flow` | Writing shallow (1-3 case) conditional logic. |
| `go-testing-conventions` | Any `_test.go` file — table-driven shape, `t.Context()`, error-path assertions. |
| `error-handling-conventions` | Any code returning/wrapping/logging an `error` — sentinel errors, `%w` wrapping, `errs.FieldErrors`. |
| `api-handler-conventions` | App-layer HTTP handlers — decode/validate/call/respond shape, status mapping, middleware. |
| `db-migration-conventions` | `business/sdk/migrate/sql` — naming, reversibility, locking, schema-mirrors-invariants. |
| `logging-conventions` | Any `log.*` call — structured fields, levels, what not to log. |
| `makefile-conventions` | Bootstrapping a new project — the make targets CI/PR-templates/review skills assume exist. |

## Review skill (read-only, advisory)

| Skill | Use for |
|---|---|
| `review-pr` | PR/diff review across all six lenses (Service Steward, Error Tripwire, Doc Drift Check, Harness Map, Boundary Keeper, Straight-Line Pass). See `review-pr/SKILL.md` for routing and the G0-G4 gate scale. |

## How the two categories relate

The review lenses *judge* a diff against the same standards the authoring
skills *teach* — e.g. `error-tripwire` (review) and `error-handling-conventions`
(authoring) cover the same ground from opposite directions. Writing to the
authoring skill up front should mean the review lens comes back clean.

## Known gap

`review-pr/SKILL.md` references a `validate-pr-title` skill ("applies whenever
the PR title is in scope, regardless of language") but no corresponding
`SKILL.md`/`reference/validate-pr-title.md` exists yet in this directory.
Either add it (a short skill checking PR titles against the conventional-commit
patterns already enforced in `.github/workflows/pr-template-check.yml`) or
remove the reference from `review-pr/SKILL.md` — don't leave it dangling.

## Adding a new skill to this boilerplate

1. `mkdir .agents/skills/<skill-name>`, add `SKILL.md` with `name` +
   `description` frontmatter (description = the trigger condition, written so
   an agent can match it against an arbitrary task description).
2. If it's a review lens, add it under `review-pr/reference/` instead and wire
   it into `review-pr/SKILL.md`'s lens catalog + routing table.
3. Add a row to the table above.