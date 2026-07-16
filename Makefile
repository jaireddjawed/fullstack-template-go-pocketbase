.PHONY: dev serve build test migrate-up migrate-down migration seed types superuser generator

# --- Running ---------------------------------------------------------------

dev: ## Run the backend with auto-applied migrations (http://127.0.0.1:8090)
	go run ./cmd/app serve

serve: dev

build: ## Build a production binary
	go build -o app ./cmd/app

generator: ## Build the project generator TUI (./fullstack-template init)
	go build -o fullstack-template ./cmd/fullstack-template

# --- Database --------------------------------------------------------------

migrate-up: ## Apply all pending migrations
	go run ./cmd/app migrate up

migrate-down: ## Revert the last applied migration
	go run ./cmd/app migrate down 1

migration: ## Create a new blank migration: make migration name=add_comments
	go run ./cmd/app migrate create "$(name)"

seed: ## Seed the database with development data (idempotent)
	go run ./cmd/app seed

superuser: ## Create a dashboard superuser: make superuser email=you@example.com pass=changeme123
	go run ./cmd/app superuser upsert "$(email)" "$(pass)"

# --- Types -----------------------------------------------------------------

types: migrate-up ## Regenerate TypeScript types from Go DTOs + PocketBase collections
	go run github.com/gzuidhof/tygo@latest generate
	npx --yes pocketbase-typegen --db pb_data/data.db --out shared/pocketbase.gen.ts

# --- Tests -----------------------------------------------------------------

test: ## Run backend tests
	go test ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-14s %s\n", $$1, $$2}'
