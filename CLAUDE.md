# CLAUDE.md

Project instructions for Claude Code (and other agents). The authoritative
rule set lives in [AGENTS.md](AGENTS.md) — read it first; this file only
adds tool-specific pointers.

- Follow every rule in `AGENTS.md` (persona "Dave", planning-first workflow,
  no automatic git commits/staging — hand commands to the user instead).
- Skills live in `.agents/skills/`. Relevant skill triggers are listed in
  `.agents/skills/README.md`. Always apply `use-modern-go` when writing Go.
- Before touching `app/*`, `business/domain/*`, or `.../stores/*db` files,
  read `.agents/skills/layered-architecture-types/SKILL.md`.
- For review requests, use the `review-pr` skill (six lenses, G0–G4 gate scale).
- Toolchain: Go (version pinned in `go.mod`). Use `make help` to list targets;
  `make check` is the full local CI equivalent. Only run `make fmt lint` after
  changing `.go` files if the user confirms.
