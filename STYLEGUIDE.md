# Coding Standards (Go)

- Go toolchain pinned in `go.mod`; prefer modern features per the
  `use-modern-go` skill (loop-range, `t.Context()`, errors.Join, etc.).
- `gofmt -s` + `goimports` with local prefix `github.com/mrbardia72/labdon`
  (`make fmt`). Imports grouped: std / external / local.
- Errors: sentinel errors + `%w` wrapping; validation via `errs.FieldErrors`;
  never swallow an error. Follow `error-handling-conventions`.
- Logging: structured (zap-style fields), contextual levels, never log secrets.
  Follow `logging-conventions`.
- Naming: mixedCaps, package names short/lowercase, no stutter
  (`business.Order`, not `business.BusinessOrder`).
- Branching: default-first assignment, naked switches over if/else ladders —
  follow `branching-logic-flow`.
- Tests: table-driven, subtests, error-path assertions first. Follow
  `go-testing-conventions`. Aim for meaningful coverage, not a number.
- Dependency injection over globals; constructors return `(*T, error)`.
- Keep functions small and linear ("straight-line" review lens).
