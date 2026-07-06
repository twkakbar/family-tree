# ============================================================
# Family Tree – developer commands
# Requires GNU make (pre-installed on Linux/macOS, `scoop install make`
# or use the raw commands from README.md on Windows).
# ============================================================

.DEFAULT_GOAL := help

## -------- Docker stacks ------------------------------------

.PHONY: dev
dev: ## Start the development stack (requires .env with COMPOSE_FILE set)
	docker compose up --build

.PHONY: dev-down
dev-down: ## Stop the development stack
	docker compose down

.PHONY: prod
prod: ## Start the production stack (VM)
	docker compose -f docker-compose.yml up -d --build

.PHONY: prod-down
prod-down: ## Stop the production stack
	docker compose -f docker-compose.yml down

.PHONY: logs
logs: ## Tail logs of all services
	docker compose logs -f --tail=100

.PHONY: ps
ps: ## Show running services
	docker compose ps

## -------- Backend -------------------------------------------

.PHONY: backend-run
backend-run: ## Run the API natively (needs local Postgres + backend/.env)
	cd backend && go run ./cmd/api

.PHONY: backend-test
backend-test: ## Run backend tests with race detector
	cd backend && go test -race -cover ./...

.PHONY: backend-lint
backend-lint: ## Vet + gofmt check
	cd backend && go vet ./... && gofmt -l .

.PHONY: backend-tidy
backend-tidy: ## Sync go.mod / go.sum
	cd backend && go mod tidy

## -------- Frontend -------------------------------------------

.PHONY: frontend-install
frontend-install: ## Install frontend dependencies
	cd frontend && npm install

.PHONY: frontend-dev
frontend-dev: ## Run Next.js dev server natively
	cd frontend && npm run dev

.PHONY: frontend-lint
frontend-lint: ## ESLint + TypeScript check
	cd frontend && npm run lint && npm run typecheck

.PHONY: frontend-build
frontend-build: ## Production build of the frontend
	cd frontend && npm run build

## -------- Database migrations --------------------------------

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations (reads backend/.env)
	cd backend && go run ./cmd/migrate up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	cd backend && go run ./cmd/migrate down 1

.PHONY: migrate-version
migrate-version: ## Show current migration version
	cd backend && go run ./cmd/migrate version

.PHONY: migrate-new
migrate-new: ## Create a new migration: make migrate-new name=create_x
	cd backend && go run ./cmd/migrate new $(name)

## -------- Help ------------------------------------------------

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
