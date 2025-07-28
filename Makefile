# TickTick MCP Server Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Build parameters
BINARY_NAME=ticktick-mcp
BINARY_CLEAN_NAME=ticktick-mcp-clean
BUILD_DIR=bin

# Version
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT := $(shell git rev-parse --short HEAD)

# Linker flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

.PHONY: help build build-clean run run-clean test test-coverage lint clean deps dev-setup

# Default target
all: build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the original version
	@echo "Building original version..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go

build-clean: ## Build the Clean Architecture version
	@echo "Building Clean Architecture version..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CLEAN_NAME) ./cmd/main_clean.go

build-all: build build-clean ## Build both versions

# Run targets
run: ## Run the original version
	@echo "Running original version..."
	$(GORUN) ./main.go

run-clean: ## Run the Clean Architecture version
	@echo "Running Clean Architecture version..."
	$(GORUN) ./cmd/main_clean.go

# Development targets
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	$(GOMOD) tidy
	$(GOMOD) download
	@echo "Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/vektra/mockery/v2@latest

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Test targets
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-domain: ## Test domain layer only
	@echo "Testing domain layer..."
	$(GOTEST) -v ./domain/...

test-usecases: ## Test use cases layer only
	@echo "Testing use cases layer..."
	$(GOTEST) -v ./usecases/...

test-adapters: ## Test adapters layer only
	@echo "Testing adapters layer..."
	$(GOTEST) -v ./adapters/...

# Quality targets
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Mock generation
mocks: ## Generate mocks for testing
	@echo "Generating mocks..."
	mockery --dir=domain/repositories --all --output=test/mocks/repositories
	mockery --dir=domain/services --all --output=test/mocks/services

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCMD) clean
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

clean-cache: ## Clean go cache
	@echo "Cleaning go cache..."
	$(GOCMD) clean -cache

# Docker targets (optional)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ticktick-mcp:$(VERSION) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -it ticktick-mcp:$(VERSION)

# Documentation targets
docs: ## Generate documentation
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./domain/...
	$(GOCMD) doc -all ./usecases/...

# Environment targets
env-example: ## Create example environment file
	@echo "Creating example environment file..."
	@echo "# TickTick API Configuration" > .env.example
	@echo "TICKTICK_CLIENT_ID=your_client_id_here" >> .env.example
	@echo "TICKTICK_CLIENT_SECRET=your_client_secret_here" >> .env.example
	@echo "TICKTICK_BASE_URL=https://api.dida365.com/open/v1" >> .env.example
	@echo "TICKTICK_AUTH_URL=https://dida365.com/oauth/authorize" >> .env.example
	@echo "TICKTICK_TOKEN_URL=https://dida365.com/oauth/token" >> .env.example
	@echo "TICKTICK_REDIRECT_URL=http://localhost:8080/callback" >> .env.example
	@echo "" >> .env.example
	@echo "# Authentication Configuration" >> .env.example
	@echo "TICKTICK_TOKEN_FILE=/tmp/ticktick_auth.json" >> .env.example
	@echo "TICKTICK_AUTH_PORT=8080" >> .env.example
	@echo "" >> .env.example
	@echo "# Server Configuration" >> .env.example
	@echo "SERVER_TIMEOUT=30" >> .env.example
	@echo "" >> .env.example
	@echo "# Storage Configuration" >> .env.example
	@echo "DATA_DIR=/tmp/ticktick-mcp" >> .env.example
	@echo "" >> .env.example
	@echo "# Logging Configuration" >> .env.example
	@echo "LOG_LEVEL=info" >> .env.example
	@echo "LOG_FORMAT=json" >> .env.example
	@echo "Example environment file created: .env.example"

# Comparison targets
compare: build build-clean ## Compare binary sizes
	@echo "Comparing binary sizes..."
	@ls -lh $(BUILD_DIR)/
	@echo ""
	@echo "Original version:"
	@du -h $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Clean Architecture version:"
	@du -h $(BUILD_DIR)/$(BINARY_CLEAN_NAME)

# Installation targets
install: build-clean ## Install the Clean Architecture version
	@echo "Installing Clean Architecture version..."
	cp $(BUILD_DIR)/$(BINARY_CLEAN_NAME) $$GOPATH/bin/$(BINARY_CLEAN_NAME)

# Quick development workflow
dev: deps fmt vet test-coverage lint ## Full development workflow

# CI/CD targets
ci: deps fmt vet test lint ## CI pipeline
	@echo "CI pipeline completed successfully"

# Architecture validation
arch-check: ## Check architecture compliance
	@echo "Checking architecture compliance..."
	@echo "Checking import dependencies..."
	@! grep -r "internal/client\|internal/server\|internal/auth" domain/ usecases/ || (echo "ERROR: Use cases or domain importing internal packages!" && exit 1)
	@! grep -r "adapters\|infrastructure\|interfaces" domain/ || (echo "ERROR: Domain importing outer layers!" && exit 1)
	@! grep -r "adapters\|infrastructure\|interfaces" usecases/ || (echo "ERROR: Use cases importing outer layers!" && exit 1)
	@echo "Architecture compliance check passed ✓"

# Performance testing
benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Security scanning
security: ## Run security scanner
	@echo "Running security scanner..."
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; $(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; }
	gosec ./...

.DEFAULT_GOAL := help