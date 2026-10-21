# DXP AI Gateway — provisioning worker
# Usage: make help

SHELL := /bin/bash
.DEFAULT_GOAL := help

DB_URL ?= postgres://dxp:dxp@localhost:5432/dxp_ai_gateway?sslmode=disable
GO ?= go

export DATABASE_URL := $(DB_URL)

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
run: ## Run worker (uses DATABASE_URL)
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
