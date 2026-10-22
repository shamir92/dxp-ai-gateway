# DXP AI Gateway — provisioning worker
# Usage: make help

SHELL := /bin/bash
.DEFAULT_GOAL := help

GO ?= go

# Configuration comes from .env (see .env.example). Nothing is exported here on
# purpose: an exported variable outranks the file, because config.LoadDotEnv only
# fills in keys that are not already in the environment.

.PHONY: help
help: ## Show targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Download Go modules
	$(GO) mod download

.PHONY: build
build: ## Compile worker binary
	$(GO) build -o bin/worker .

.PHONY: run
run: ## Run worker (reads configuration from .env)
	$(GO) run .

.PHONY: test
test: ## Unit tests
	$(GO) test ./...

.PHONY: vet
vet: ## go vet
	$(GO) vet ./...

.PHONY: fmt
fmt: ## gofmt write
	gofmt -w .

.PHONY: tidy
tidy: ## go mod tidy
	$(GO) mod tidy

.PHONY: docker
docker: ## Build the worker container image
	docker build -t dxp-ai-gateway/worker:local .

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin
