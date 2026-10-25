# Architecture Decision Records (ADRs)

Lightweight decision log. Newest entries at the top. Keep each entry short:
**Context → Decision → Consequences.**

---

## ADR-0004 — CI runs through the Makefile
- **Date:** 2026-10-25
- **Status:** Accepted
- **Context:** The two GitHub workflows previously duplicated raw `go test`/`gofmt`/build commands with drifted paths and Go versions, so local and CI behavior diverged.
- **Decision:** All pipelines call `make ci` / `make check` (the same targets developers run locally). Raw command duplication in workflows is not allowed.
- **Consequences:** One source of truth; fixing CI means fixing the Makefile.

## ADR-0003 — Go version policy
- **Date:** 2026-10-25
- **Status:** Accepted
- **Context:** Workflows referenced 1.23–1.27 inconsistently.
- **Decision:** Single supported toolchain: the version in `go.mod` (`go 1.24`). Bump it deliberately via an ADR, never ad-hoc per workflow.
- **Consequences:** Predictable builds; matrix removed (older Go releases are EOL'd by the Go project's own policy).

## ADR-0002 — Layered architecture (app / business / foundation)
- **Date:** 2026-10-25
- **Status:** Accepted
- **Context:** Docs conflicted (generic API→Service→Repository vs the app/business/stores style described in AGENTS.md and the skills).
- **Decision:** Follow the layered model taught by `.agents/skills/layered-architecture-types`: primitives at the edges (app JSON, DB rows), strong types only in `business/types/*`, named converters (`toBus<Type>`, `fromBus<Type>Response`, `toDB<Type>`) at every crossing.
- **Consequences:** ARCHITECTURE.md and STYLEGUIDE.md rewritten to match; review lenses enforce it.

## ADR-0001 — Language: Go
- **Date:** 2026-10-25
- **Status:** Accepted
- **Context:** STYLEGUIDE.md mentioned TypeScript strict mode while the entire toolchain (Makefile, workflows, skills) is Go-centric.
- **Decision:** This is a Go-only project. TypeScript references were stale copies from another repo and are removed.
- **Consequences:** Docs unified on Go + modern toolchain features.
