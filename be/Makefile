.PHONY: help dev run test build tidy fmt swagger install-swagger install-migrate migrate-up migrate-down migrate-force migrate-version migrate-create

MIGRATE := "$(shell go env GOPATH)/bin/migrate"
seed ?= false

help: ## Show available commands
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sed 's/:.*## / - /'

dev: ## Run live reload server with Air (usage: make dev seed=true)
	SEED_ON_STARTUP=$(seed) go tool air

run: ## Run the app once (usage: make run seed=true)
	SEED_ON_STARTUP=$(seed) go run .

test: ## Run all tests
	go test ./...

build: ## Build binary to ./bin/api
	mkdir -p bin && go build -o bin/api .

tidy: ## Tidy Go modules
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

install-swagger: ## Install Swagger dependencies and swag CLI
	go get github.com/swaggo/gin-swagger github.com/swaggo/files github.com/swaggo/swag@latest
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: ## Generate Swagger docs in ./docs
	"$(shell go env GOPATH)/bin/swag" init -g main.go -o docs

install-migrate: ## Install golang-migrate CLI (with postgres driver)
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Run all up migrations against database in .env
	@set -a; . ./.env; set +a; $(MIGRATE) -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" up

migrate-down: ## Roll back one migration against database in .env
	@set -a; . ./.env; set +a; $(MIGRATE) -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" down 1

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@set -a; . ./.env; set +a; $(MIGRATE) -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" force $(version)

migrate-version: ## Show current migration version
	@set -a; . ./.env; set +a; $(MIGRATE) -path ./migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" version

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users)
	$(MIGRATE) create -ext sql -dir ./migrations -seq $(name)
