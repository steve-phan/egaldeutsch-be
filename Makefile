.PHONY: help build run test clean docker-build docker-run docker-stop migrate-up migrate-down

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Go commands
dev: ## Run the application in development mode
	go run ./cmd/server

build: ## Build the application
	go build -o bin/server ./cmd/server

run: ## Run the application
	go run ./cmd/server

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -rf bin/

# Docker commands
docker-build: ## Build Docker image
	docker build -t egaldeutsch-be .

docker-run: ## Run with Docker Compose
	docker-compose up --build

docker-stop: ## Stop Docker Compose services
	docker-compose down

# Database migration commands (requires migrate tool)
migrate-up: ## Run database migrations up
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/egaldeutsch-go?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/egaldeutsch-go?sslmode=disable" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	migrate create -ext sql -dir migrations -seq $(name)

# Development setup
setup: ## Setup development environment
	go mod tidy
	cp .env.example .env

# Linting and formatting
fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: fmt vet ## Run formatting and vetting