---
name: makefile-conventions
description: Defines the standard Makefile targets this repo's CI, PR templates, and review skills assume exist (make test, make test-only, make lint, ...). Use when creating a new Go project from this boilerplate, or when adding/renaming a make target.
---

# Makefile Conventions

Every CI workflow, PR template, and the `harness-map` review lens in this repo
refer to make targets by name (`make test`, `make lint`, `make test-only`).
None of those references work unless a Makefile actually defines them. This
skill is the contract: when bootstrapping a new Go project from this
boilerplate, create a `Makefile` exposing at least these targets, with these
exact names and this behavior.

## Required targets

| Target | Behavior |
|---|---|
| `test` | Full test suite, the way CI runs it: `go test -race -coverprofile=coverage.out ./...` at minimum. Should fail (non-zero exit) on any test failure — don't swallow exit codes. |
| `test-only` | Same as `test` but without coverage/race overhead, for fast local iteration: `go test ./...`. `harness-map`'s "Suggested verification" section assumes this exists and is fast. |
| `lint` | Runs `golangci-lint run --config=.golangci.yml` (matching the version pinned in `.github/workflows/go-lint.yml`'s `GOLANGCI_LINT_VERSION`) plus `gofmt -s -l .` / `go vet ./...` if golangci-lint's config doesn't already cover them. |
| `build` | Builds the primary binary/binaries for the current `GOOS`/`GOARCH`, mirroring (a subset of) the cross-compile matrix in `ci-cd.yml`'s "Build application" step. |
| `tidy` | `go mod tidy` — the same command `go-basic-checks` diffs against in CI; running it locally before pushing avoids a CI failure on an out-of-date `go.sum`. |
| `fmt` | `gofmt -s -w .` and `goimports -w .` — the writable counterpart to the `go fmt check` / `goimports check` CI steps, which only read. |

## Naming and conventions

- Targets are lowercase, hyphenated for multi-word names (`test-only`, not
  `testOnly` or `test_only`).
- Every target marks itself `.PHONY` — none of these produce a file with the
  same name as the target.
- `make` with no target prints help (a `help` target listing available
  commands), rather than running an arbitrary default.
- Keep the Makefile itself thin: it should shell out to `go`, `golangci-lint`,
  etc., not contain real logic. Put real logic in `bot/cmd/tooling` or a script
  if it grows beyond a couple of lines per target.

## Why this matters for this boilerplate

Several other pieces in `.agents/skills` and `.github` assume these targets
exist and behave this way:

- `go-lint.yml` and `ci-cd.yml` could be pointed at `make lint` / `make test`
  instead of inlining the raw commands — keeping the Makefile and CI in sync
  is then a matter of editing one place.
- `harness-map.md`'s "Suggested verification" section literally prints
  `make test-only` and `go test ./path/...` as the two recommended commands —
  if `test-only` doesn't exist, that guidance is dead on arrival for this repo.
- The PR templates' checklists (`simple.md`, `bug_fix.md`, the default
  template) all have a literal `make test` / `make lint` checkbox — a new
  project using these templates without a matching Makefile makes that
  checklist item unverifiable.

## Minimal starter Makefile

```makefile
.PHONY: help test test-only lint build tidy fmt

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | sed 's/:.*##/:/'

test: ## Full test suite with race detector and coverage
	go test -race -coverprofile=coverage.out ./...

test-only: ## Fast local test run, no race/coverage
	go test ./...

lint: ## Run golangci-lint plus fmt/vet checks
	golangci-lint run --config=.golangci.yml
	gofmt -s -l .
	go vet ./...

build: ## Build for the current platform
	go build -o bin/ ./...

tidy: ## Tidy go.mod/go.sum
	go mod tidy

fmt: ## Auto-format and fix imports
	gofmt -s -w .
	goimports -w .
```

Adjust binary paths/targets in `build` to match the project's actual `cmd/`
layout (see `ci-cd.yml`'s build step for this repo's specific binary names).