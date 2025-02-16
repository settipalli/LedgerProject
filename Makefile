# Build configuration
BINARY_NAME=ledger
MAIN_PACKAGE=main.go
GO=go

# Directories
BUILD_DIR=build
COVERAGE_DIR=coverage

# Environment variables
GOFLAGS=-mod=vendor
ENVIRONMENT?=development

# Git information
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Linting tools
GOLANGCI_LINT=golangci-lint

.PHONY: all
all: clean deps test lint build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="-X main.Version=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	APP_ENV=$(ENVIRONMENT) ./$(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod vendor

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@mkdir -p $(COVERAGE_DIR)
	APP_ENV=test $(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Lint the code
.PHONY: lint
lint:
	@echo "Linting code..."
	$(GOLANGCI_LINT) run

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	go clean -testcache

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Security check
.PHONY: security
security:
	@echo "Running security checks..."
	gosec ./...

# Development helpers
.PHONY: dev
dev: deps fmt lint test build

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, download dependencies, test, lint, and build"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  deps         - Download dependencies"
	@echo "  test         - Run tests with coverage"
	@echo "  bench        - Run benchmarks"
	@echo "  lint         - Lint the code"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  security     - Run security checks"
	@echo "  dev          - Development workflow (deps, fmt, lint, test, build)"
	@echo "  help         - Show this help message"

# Version information
.PHONY: version
version:
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
