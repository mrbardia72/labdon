SHELL := /bin/bash
.DEFAULT_GOAL := help

# ── Configuration ────────────────────────────────────────────────────────────
# Override any of these on the command line, e.g.:
#   make build BINARY_NAME=myservice CMD_PATH=./cmd/server
GO              ?= go
GOLANGCI_LINT   ?= golangci-lint
BINARY_NAME     ?= app
CMD_PATH        ?= ./cmd/$(BINARY_NAME)
BIN_DIR         ?= bin
COVERAGE_OUT    ?= coverage.out
GOLANGCI_CONFIG ?= .golangci.yml

# Cross-compile targets used by `make build-all` (GOOS/GOARCH pairs)
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: help
help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ── Dependencies ─────────────────────────────────────────────────────────────

.PHONY: deps
deps: ## Download go modules
	$(GO) mod download

.PHONY: tidy
tidy: ## Tidy go.mod/go.sum
	$(GO) mod tidy

.PHONY: verify
verify: ## Verify module checksums
	$(GO) mod verify

.PHONY: tidy-check
tidy-check: ## Fail if go.mod/go.sum are not tidy (used in CI)
	$(GO) mod tidy
	@git diff --exit-code go.mod go.sum || \
		(echo "go.mod/go.sum are not tidy — run 'make tidy' and commit the result" && exit 1)

# ── Formatting ───────────────────────────────────────────────────────────────

.PHONY: fmt
fmt: ## Auto-format code and fix imports
	gofmt -s -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w . || \
		echo "goimports not installed — run 'make tools' first"

.PHONY: fmt-check
fmt-check: ## Fail if any file is not gofmt-formatted (used in CI)
	@unformatted=$$(gofmt -s -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$unformatted"; \
		echo "Run 'make fmt' to fix."; \
		exit 1; \
	fi

# ── Vet & lint ───────────────────────────────────────────────────────────────

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint (requires .golangci.yml)
	$(GOLANGCI_LINT) run --timeout=10m --config=$(GOLANGCI_CONFIG)

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with autofix
	$(GOLANGCI_LINT) run --timeout=10m --config=$(GOLANGCI_CONFIG) --fix

# ── Tests ────────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Full test suite with race detector and coverage (used in CI)
	$(GO) test -v -race -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...

.PHONY: test-only
test-only: ## Fast local test run, no race detector or coverage
	$(GO) test ./...

.PHONY: test-bench
test-bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

.PHONY: coverage
coverage: test ## Generate and open an HTML coverage report
	$(GO) tool cover -html=$(COVERAGE_OUT) -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: coverage-func
coverage-func: test ## Print per-function coverage summary
	$(GO) tool cover -func=$(COVERAGE_OUT)

# ── Security & static analysis ───────────────────────────────────────────────

.PHONY: security
security: ## Run gosec and govulncheck
	@command -v gosec >/dev/null 2>&1 || $(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec ./...
	@command -v govulncheck >/dev/null 2>&1 || $(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: staticcheck
staticcheck: ## Run staticcheck
	@command -v staticcheck >/dev/null 2>&1 || $(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# ── Build ────────────────────────────────────────────────────────────────────

.PHONY: build
build: ## Build the binary for the current platform
	$(GO) build -v -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

.PHONY: build-all
build-all: ## Cross-compile for all platforms in $(PLATFORMS)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		ext=""; if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		out="$(BIN_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext"; \
		echo "Building $$out"; \
		GOOS=$$os GOARCH=$$arch $(GO) build -v -o "$$out" $(CMD_PATH); \
	done

.PHONY: install
install: ## Install the binary into $GOPATH/bin
	$(GO) install $(CMD_PATH)

# ── Run ──────────────────────────────────────────────────────────────────────

.PHONY: run
run: ## Run the service locally
	$(GO) run $(CMD_PATH)

# ── Tooling ──────────────────────────────────────────────────────────────────

.PHONY: tools
tools: ## Install the dev tools used by the targets above
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "Note: golangci-lint is intentionally not auto-installed here —"
	@echo "install the pinned version from https://golangci-lint.run/welcome/install/"
	@echo "to match the version your CI workflow uses."

# ── Cleanup ──────────────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Remove build artifacts and coverage output
	rm -rf $(BIN_DIR) $(COVERAGE_OUT) coverage.html

# ── Composite targets ────────────────────────────────────────────────────────

.PHONY: check
check: tidy-check fmt-check vet lint test ## Run everything CI runs, in one shot

.PHONY: ci
ci: deps verify check ## Full CI-equivalent pipeline, for local dry-runs before pushing