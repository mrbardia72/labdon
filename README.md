# ── labdon ────────────────────────────────────────────────────────────────────

[![CI](https://github.com/mrbardia72/labdon/actions/workflows/ci.yml/badge.svg)](https://github.com/mrbardia72/labdon/actions/workflows/ci.yml)
[![Release](https://github.com/mrbardia72/labdon/actions/workflows/release.yml/badge.svg)](https://github.com/mrbardia72/labdon/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrbardia72/labdon)](https://goreportcard.com/report/github.com/mrbardia72/labdon)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Go service scaffold with a strict layered architecture (`app` → `business` →
`foundation`) and a full set of AI-agent skills under `.agents/skills/` that
teach and enforce the repo's conventions. The current code is a working
**health-check server** intended as the foundation to build on.

## Features

- 🏛️ **Layered architecture** — primitives at the edges (API JSON, DB rows),
  strong types only in the business layer, named converters at every crossing.
- 🤖 **Agent-ready** — `AGENTS.md`, `CLAUDE.md`, `.cursorrules`,
  Copilot instructions, and 10+ review/authoring skills for AI coding agents.
- 🧰 **One Makefile source of truth** — CI runs exactly what you run locally
  (`make ci` / `make check`).
- 🔒 **Security baked in** — gosec, govulncheck, Trivy scanning; graceful
  shutdown and request timeouts in the server.
- ✅ **Quality gates** — golangci-lint, race-detector tests, coverage,
  conventional-commit PR checks, PR-size analysis.

## Requirements

- Go **1.24+** (pinned in [`go.mod`](go.mod))
- GNU Make
- Optional dev tools: `golangci-lint`, `goimports`, `gosec`, `govulncheck`,
  `staticcheck` → install with `make tools`

## Quickstart

```bash
git clone https://github.com/mrbardia72/labdon.git
cd labdon
make deps     # download modules
make run      # start the service on :3000 (default)
curl localhost:3000/v1/liveness
curl localhost:3000/v1/readiness
```

Configure via environment variables or `config.dev.json`:

```bash
LABDON_WEB_HTTPREADTIMEOUT=5s LABDON_WEB_HTTPWRITE_TIMEOUT=10s LABDON_WEB_HTTPPORT=8080 make run
```

## Project layout

```
cmd/sites            # main entrypoint (flag parsing, DI wiring, shutdown)
app                  # HTTP layer: handlers, middleware (primitives only)
business             # domain logic + strong types (no cross-domain imports)
foundation           # shared infra: logger, trace, database, buildinfo
business/sdk/errs    # typed error helpers (FieldErrors, New)
```

Full rules: [ARCHITECTURE.md](ARCHITECTURE.md) · coding standards:
[STYLEGUIDE.md](STYLEGUIDE.md) · decisions: [DECISIONS.md](DECISIONS.md)

## Development

```bash
make help         # list all targets
make fmt lint     # format + lint after touching .go files
make test-only    # fast tests
make check        # everything CI runs (tidy, fmt, vet, lint, test)
make build-all    # cross-compile linux/darwin/windows × amd64/arm64
make security     # gosec + govulncheck
```

## Releases

Push a tag to publish binaries via GitHub Release:

```bash
git tag v0.1.0 && git push origin v0.1.0
```

## Contributing

Read [AGENTS.md](AGENTS.md) first (it defines the repo's agent contract), then
open a PR using one of the templates in `.github/PULL_REQUEST_TEMPLATE/`.
PR titles must follow conventional commits — validated in CI.

## License

[MIT](LICENSE)
