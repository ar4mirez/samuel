# AICoF CLI Makefile

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X github.com/ar4mirez/aicof/internal/commands.Version=$(VERSION) \
	-X github.com/ar4mirez/aicof/internal/commands.Commit=$(COMMIT) \
	-X github.com/ar4mirez/aicof/internal/commands.BuildDate=$(BUILD_DATE)"

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOLINT := golangci-lint

# Binary name
BINARY_NAME := aicof
BINARY_PATH := ./bin/$(BINARY_NAME)

# Main package
MAIN_PACKAGE := ./cmd/aicof

.PHONY: all build clean test lint fmt deps help install uninstall

## Default target
all: deps lint test build

## Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p ./bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PACKAGE)
	@echo "Built: $(BINARY_PATH)"

## Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p ./bin
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ./bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ./bin/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ./bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Built all platforms in ./bin/"

## Install locally
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@cp $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)
	@echo "Installed. Run 'aicof version' to verify."

## Uninstall
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstalled."

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ./bin
	@$(GOCMD) clean
	@echo "Clean complete."

## Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -cover ./...

## Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## Run linter
lint:
	@echo "Running linter..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w -s .
	@echo "Format complete."

## Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies ready."

## Run the CLI (dev mode)
run:
	@$(GOCMD) run $(MAIN_PACKAGE) $(ARGS)

## Show version info
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

## Release with goreleaser (dry run)
release-dry:
	@echo "Running goreleaser (dry run)..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest"; \
	fi

## Release with goreleaser
release:
	@echo "Running goreleaser..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --clean; \
	else \
		echo "goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest"; \
	fi

## Help
help:
	@echo "AICoF CLI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all           Run deps, lint, test, and build (default)"
	@echo "  build         Build the binary"
	@echo "  build-all     Build for all platforms"
	@echo "  install       Install to /usr/local/bin"
	@echo "  uninstall     Remove from /usr/local/bin"
	@echo "  clean         Remove build artifacts"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  lint          Run linter"
	@echo "  fmt           Format code"
	@echo "  deps          Download dependencies"
	@echo "  run           Run CLI (use ARGS=\"...\" for arguments)"
	@echo "  version       Show version info"
	@echo "  release-dry   Test release with goreleaser"
	@echo "  release       Release with goreleaser"
	@echo "  help          Show this help"
