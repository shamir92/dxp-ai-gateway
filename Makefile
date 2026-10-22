# DXP AI Gateway — root
# Usage: make help

SHELL := /bin/bash
.DEFAULT_GOAL := help

ROOT := $(CURDIR)
BACKEND := $(ROOT)/lp/backend
FRONTEND := $(ROOT)/lp/frontend
WORKER := $(ROOT)/lp/worker

# DB_URL ?= postgres://dxp:dxp@localhost:5432/dxp_ai_gateway?sslmode=disable
HTTP_ADDR ?= :8080
PUBLIC_WEB_URL ?= http://localhost:3000
GO ?= go

export HTTP_ADDR
export PUBLIC_WEB_URL

.PHONY: help
help: ## Show targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# --- database ---

.PHONY: db-up
db-up: ## Start local Postgres (docker compose in lp/backend)
	$(MAKE) -C $(BACKEND) db-up

.PHONY: db-down
db-down: ## Stop Postgres
	$(MAKE) -C $(BACKEND) db-down

.PHONY: db-reset
db-reset: ## Wipe Postgres volume and recreate
	$(MAKE) -C $(BACKEND) db-reset

.PHONY: psql
psql: ## Open psql in the app database
	$(MAKE) -C $(BACKEND) psql

# --- backend ---

.PHONY: deps
deps: ## Download Go modules
	$(MAKE) -C $(BACKEND) deps

.PHONY: build
build: ## Compile backend server
	$(MAKE) -C $(BACKEND) build

.PHONY: run
run: ## Run backend API server
	$(MAKE) -C $(BACKEND) run

.PHONY: worker
worker: ## Compile provisioning worker binary
	$(MAKE) -C $(WORKER) build

.PHONY: run-worker
run-worker: ## Run provisioning worker
	$(MAKE) -C $(WORKER) run

.PHONY: worker-docker
worker-docker: ## Build provisioning worker container image
	$(MAKE) -C $(WORKER) docker

.PHONY: test
test: ## Backend unit tests
	$(MAKE) -C $(BACKEND) test

.PHONY: test-verbose
test-verbose: ## Backend unit tests (verbose)
	$(MAKE) -C $(BACKEND) test-verbose

.PHONY: vet
vet: ## go vet backend
	$(MAKE) -C $(BACKEND) vet

.PHONY: fmt
fmt: ## gofmt backend
	$(MAKE) -C $(BACKEND) fmt

.PHONY: tidy
tidy: ## go mod tidy backend
	$(MAKE) -C $(BACKEND) tidy

.PHONY: backend
backend: build ## Alias for backend build

# --- frontend ---

.PHONY: frontend
frontend: ## Build Next.js frontend
	$(MAKE) -C $(BACKEND) frontend

.PHONY: frontend-dev
frontend-dev: ## Run Next.js dev server
	$(MAKE) -C $(BACKEND) frontend-dev

.PHONY: frontend-install
frontend-install: ## npm install frontend
	cd $(FRONTEND) && npm install

# --- workflows ---

.PHONY: dev
dev: ## Start Postgres + backend (console mail)
	$(MAKE) -C $(BACKEND) dev

.PHONY: dev-all
dev-all: ## Postgres + backend + frontend (two terminals after)
	@echo "1) make db-up"
	@echo "2) make run          # terminal A"
	@echo "3) make frontend-dev # terminal B"
	$(MAKE) db-up
	$(MAKE) run

.PHONY: smoke
smoke: ## curl backend /healthz
	$(MAKE) -C $(BACKEND) smoke

.PHONY: test-all
test-all: test ## Backend tests + frontend typecheck/build
	$(MAKE) frontend

.PHONY: clean
clean: ## Remove backend build artifacts
	$(MAKE) -C $(BACKEND) clean
