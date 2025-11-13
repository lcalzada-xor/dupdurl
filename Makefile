# Makefile for dupdurl

.PHONY: all build install test test-unit test-integration test-coverage bench clean fmt lint help

# Variables
BINARY_NAME=dupdurl
MAIN_FILE=main_new.go
GO=go
GOFLAGS=-ldflags="-s -w"
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Build information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

all: build ## Build the binary

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_FILE) cmd/dupdurl/cli.go

build-all: ## Build binaries for all platforms
	@echo "Building for Linux amd64..."
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-linux-amd64 $(MAIN_FILE) cmd/dupdurl/cli.go
	@echo "Building for Linux arm64..."
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-linux-arm64 $(MAIN_FILE) cmd/dupdurl/cli.go
	@echo "Building for macOS amd64..."
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_FILE) cmd/dupdurl/cli.go
	@echo "Building for macOS arm64..."
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-darwin-arm64 $(MAIN_FILE) cmd/dupdurl/cli.go
	@echo "Building for Windows amd64..."
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE) cmd/dupdurl/cli.go

install: build ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(GOFLAGS) -ldflags="$(LDFLAGS)"

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test -v -race ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GO) test -v -race ./tests/unit/...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	$(GO) test -v -race ./tests/integration/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@$(GO) tool cover -func=$(COVERAGE_FILE) | grep total

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem -run=^$$ ./tests/benchmark/...

bench-compare: ## Run benchmarks and compare with previous results
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem -run=^$$ ./tests/benchmark/... | tee benchmark-new.txt
	@if [ -f benchmark-old.txt ]; then \
		echo "Comparing with previous results..."; \
		benchstat benchmark-old.txt benchmark-new.txt; \
	fi
	@mv benchmark-new.txt benchmark-old.txt

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GO) fmt ./...
	gofmt -s -w .

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	$(GO) mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f benchmark-*.txt
	$(GO) clean

demo: build ## Run demo with test data
	@echo "Running demo..."
	@cat tests/fixtures/test_urls.txt | ./$(BINARY_NAME) -stats

demo-fuzzy: build ## Run demo with fuzzy mode
	@echo "Running demo with fuzzy mode..."
	@cat tests/fixtures/test_urls.txt | ./$(BINARY_NAME) -fuzzy -stats

demo-json: build ## Run demo with JSON output
	@echo "Running demo with JSON output..."
	@cat tests/fixtures/test_urls.txt | ./$(BINARY_NAME) -output=json

demo-parallel: build ## Run demo with parallel processing
	@echo "Running demo with parallel processing..."
	@cat tests/fixtures/test_urls.txt | ./$(BINARY_NAME) -workers=4 -stats

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

check: test lint vet ## Run all checks (tests, linting, vetting)

ci: clean deps check build ## Run CI pipeline locally

release: clean test build-all ## Build release binaries
	@echo "Release binaries built successfully"

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
