.PHONY: dev serve backend build test migrate-up migrate-down migration seed types superuser

# --- Running ---------------------------------------------------------------

dev: ## Run the backend + Vite frontend (open http://127.0.0.1:8090)
	@if [ ! -d frontend/node_modules ]; then \
		echo "frontend/node_modules is missing. Run: cd frontend && npm install"; \
		exit 1; \
	fi; \
	trap 'kill 0' INT TERM EXIT; \
	npm --prefix frontend run dev & \
	go run . serve

backend: ## Run only the backend with auto-applied migrations
	go run . serve

serve: backend

build: ## Build a production binary
	go build -o app .

# --- Database --------------------------------------------------------------

migrate-up: ## Apply all pending migrations
	go run . migrate up

migrate-down: ## Revert the last applied migration
	go run . migrate down 1

migration: ## Create a new blank migration: make migration name=add_comments
	go run . migrate create "$(name)"

seed: ## Seed the database with development data (idempotent)
	go run . seed

superuser: ## Create a dashboard superuser: make superuser email=you@example.com pass=changeme123
	go run . superuser upsert "$(email)" "$(pass)"

# --- Types -----------------------------------------------------------------

types: migrate-up ## Regenerate TypeScript types from Go DTOs + PocketBase collections
	go run github.com/gzuidhof/tygo@latest generate
	npx --yes pocketbase-typegen --db pb_data/data.db --out shared/pocketbase.gen.ts

# --- Tests -----------------------------------------------------------------

test: ## Run backend tests
	go test ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-14s %s\n", $$1, $$2}'
