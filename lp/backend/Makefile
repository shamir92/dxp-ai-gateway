# DXP AI Gateway — backend
# Usage: make help

SHELL := /bin/bash
.DEFAULT_GOAL := help

# DB_URL ?= postgres://dxp:dxp@localhost:5432/dxp_ai_gateway?sslmode=disable
HTTP_ADDR ?= :8080
PUBLIC_WEB_URL ?= http://localhost:3000
# MAIL_BACKEND ?= console
GO ?= go

# export DATABASE_URL := $(DB_URL)
export HTTP_ADDR
export PUBLIC_WEB_URL
# export MAIL_BACKEND

.PHONY: help
help: ## Show targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: db-up
db-up: ## Start local Postgres (docker compose)
	docker compose up -d

.PHONY: db-down
db-down: ## Stop Postgres
	docker compose down

.PHONY: db-reset
db-reset: ## Wipe Postgres volume and recreate
	docker compose down -v
	docker compose up -d

.PHONY: deps
deps: ## Download Go modules
	$(GO) mod download

.PHONY: build
build: ## Compile server binary
	$(GO) build -o bin/server ./cmd/server

.PHONY: run
run: ## Run API server (uses DATABASE_URL / MAIL_BACKEND)
	$(GO) run ./cmd/server

.PHONY: dev
dev: db-up ## Start Postgres + API with console mail
	@echo "DB ready (or still starting). Ctrl+C stops the server only."
	$(MAKE) run

.PHONY: test
test: ## Unit tests
	$(GO) test ./...

.PHONY: test-verbose
test-verbose: ## Unit tests (verbose)
	$(GO) test ./... -v

.PHONY: vet
vet: ## go vet
	$(GO) vet ./...

.PHONY: fmt
fmt: ## gofmt write
	gofmt -w ./cmd ./internal

.PHONY: tidy
tidy: ## go mod tidy
	$(GO) mod tidy

.PHONY: migrate
migrate: ## Apply migrations (golang-migrate via go run path or server auto-migrate)
	$(GO) run ./cmd/server &
	@sleep 2; pkill -f 'bin/server|go run ./cmd/server' || true
	@echo "Migrations also run automatically on server start."

.PHONY: psql
psql: ## Open psql in the app database
	docker compose exec postgres psql -U dxp -d dxp_ai_gateway

.PHONY: smoke
smoke: ## curl healthz (server must be running)
	curl -sS http://localhost$(HTTP_ADDR)/healthz | head -c 200; echo

.PHONY: frontend
frontend: ## Build Next.js frontend
	cd ../frontend && npm run build

.PHONY: frontend-dev
frontend-dev: ## Run Next.js dev server
	cd ../frontend && npm run dev

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin
