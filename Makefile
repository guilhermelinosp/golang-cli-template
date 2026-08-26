# golang-cli-template — development tasks
# Targets are stable API for developers AND CI: avoid renaming/removing.
SHELL := /usr/bin/env bash
GO ?= go
APP ?= golang-cli-template
CMD_DIR := ./cmd/app

# --- Build metadata (injected, never hardcoded — see internal/build) -------
ifeq ($(strip $(VERSION)),)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
endif
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE ?= $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X github.com/guilhermelinosp/golang-cli-template/internal/build.Version=$(VERSION) \
	-X github.com/guilhermelinosp/golang-cli-template/internal/build.Commit=$(COMMIT) \
	-X github.com/guilhermelinosp/golang-cli-template/internal/build.Date=$(DATE)

.PHONY: help setup fmt vet lint test cover race build run release clean verify

help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

setup: ## Initialize template: prompts for app/module names and rewires everything
	@./scripts/setup.sh

fmt: ## Format code
	$(GO) fmt ./...

vet: ## Static analysis (stdlib)
	$(GO) vet ./...

lint: ## Lint via golangci-lint
	golangci-lint run

test: ## Run the full test suite
	$(GO) test -count=1 ./...

race: ## Run tests with race detector
	$(GO) test -race -count=1 ./...

cover: ## Coverage summary plus HTML report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -1
	$(GO) tool cover -html=coverage.out -o coverage.html

build: ## Compile bin/$(APP) with stamped metadata
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(APP) $(CMD_DIR)

run: ## Execute from source (dev builds show version=dev)
	$(GO) run $(CMD_DIR)

release: ## Snapshot release into dist/ (CI publishes tagged releases)
	goreleaser release --snapshot --clean

clean: ## Remove build artifacts
	rm -rf bin/ dist/ coverage.out coverage.html

verify: ## Architectural guardrail: Cobra stays inside internal/cli/cobra
	@violators=$$(grep -rl 'spf13/cobra\|spf13/pflag' --include='*.go' . | grep -v '^./internal/cli/cobra/' || true); \
	if [ -n "$$violators" ]; then \
		echo "BOUNDARY VIOLATION — engine types leaked outside internal/cli/cobra:"; \
		echo "$$violators"; exit 1; \
	fi; \
	echo "engine boundary OK"
