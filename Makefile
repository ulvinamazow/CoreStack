.PHONY: help build run test clean docker-build docker-up docker-down docker-logs migrate

# Variables
BINARY_NAME=corestack
GO=go
DOCKER_IMAGE=corestack
DOCKER_COMPOSE=docker compose

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	$(GO) mod download
	$(GO) mod tidy

build: ## Build the application
	$(GO) build -o $(BINARY_NAME) cmd/server/main.go

run: ## Run the application locally
	$(GO) run cmd/server/main.go

test: ## Run tests
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage report
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	$(GO) clean

fmt: ## Format code
	$(GO) fmt ./...

lint: ## Run linter
	golangci-lint run

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-up: ## Start Docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	$(DOCKER_COMPOSE) down

docker-logs: ## View Docker logs
	$(DOCKER_COMPOSE) logs -f api

docker-clean: ## Remove Docker containers and volumes
	$(DOCKER_COMPOSE) down -v

migrate: ## Run database migrations
	$(GO) run cmd/server/main.go

.DEFAULT_GOAL := help
