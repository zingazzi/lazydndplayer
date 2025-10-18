# Makefile for lazydndplayer
# D&D character manager TUI application

# Variables
BINARY_NAME=lazydndplayer
GO=go
GOFLAGS=-v
TEST_FLAGS=-v -race -timeout 30s
COVERAGE_FILE=coverage.out
BUILD_DIR=.

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m
COLOR_RED=\033[31m

.PHONY: all help build test test-verbose test-coverage clean run fmt vet lint install deps check rebuild

# Default target
all: check build

## help: Show this help message
help:
	@echo "$(COLOR_BOLD)lazydndplayer - Makefile Commands$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Building:$(COLOR_RESET)"
	@echo "  make build         - Build the application"
	@echo "  make rebuild       - Clean and rebuild"
	@echo "  make install       - Install dependencies"
	@echo ""
	@echo "$(COLOR_BLUE)Testing:$(COLOR_RESET)"
	@echo "  make test          - Run all tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-models   - Run only model tests"
	@echo ""
	@echo "$(COLOR_BLUE)Code Quality:$(COLOR_RESET)"
	@echo "  make fmt           - Format code with gofmt"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run golangci-lint (if installed)"
	@echo "  make check         - Run fmt, vet, and test"
	@echo ""
	@echo "$(COLOR_BLUE)Running:$(COLOR_RESET)"
	@echo "  make run           - Build and run the application"
	@echo "  make clean-run     - Clean character data and run fresh"
	@echo ""
	@echo "$(COLOR_BLUE)Cleanup:$(COLOR_RESET)"
	@echo "  make clean         - Remove binary and build artifacts"
	@echo "  make clean-data    - Remove character save data"
	@echo "  make clean-all     - Remove everything (binary + data + coverage)"
	@echo ""

## deps: Install/update dependencies
deps:
	@echo "$(COLOR_BLUE)Installing dependencies...$(COLOR_RESET)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies installed$(COLOR_RESET)"

## build: Build the application
build:
	@echo "$(COLOR_BLUE)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .
	@echo "$(COLOR_GREEN)✓ Build complete: $(BINARY_NAME)$(COLOR_RESET)"
	@ls -lh $(BINARY_NAME)

## rebuild: Clean and rebuild
rebuild: clean build

## test: Run all tests
test:
	@echo "$(COLOR_BLUE)Running all tests...$(COLOR_RESET)"
	@$(GO) test ./tests/... $(TEST_FLAGS)
	@echo "$(COLOR_GREEN)✓ All tests passed$(COLOR_RESET)"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(COLOR_BLUE)Running tests (verbose)...$(COLOR_RESET)"
	@$(GO) test ./tests/... -v -race

## test-models: Run only model tests
test-models:
	@echo "$(COLOR_BLUE)Running model tests...$(COLOR_RESET)"
	@$(GO) test ./tests/models -v

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	@$(GO) test ./tests/... -coverprofile=$(COVERAGE_FILE) -covermode=atomic
	@$(GO) tool cover -func=$(COVERAGE_FILE)
	@echo ""
	@echo "$(COLOR_YELLOW)To view HTML coverage report, run:$(COLOR_RESET)"
	@echo "  go tool cover -html=$(COVERAGE_FILE)"
	@echo "$(COLOR_GREEN)✓ Coverage report generated: $(COVERAGE_FILE)$(COLOR_RESET)"

## test-coverage-html: Run tests with coverage and open HTML report
test-coverage-html: test-coverage
	@echo "$(COLOR_BLUE)Opening coverage report in browser...$(COLOR_RESET)"
	@$(GO) tool cover -html=$(COVERAGE_FILE)

## fmt: Format code with gofmt
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@$(GO) fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	@$(GO) vet ./...
	@echo "$(COLOR_GREEN)✓ Vet checks passed$(COLOR_RESET)"

## lint: Run golangci-lint (if installed)
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Lint checks passed$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not installed, skipping...$(COLOR_RESET)"; \
		echo "  Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## check: Run fmt, vet, and test
check: fmt vet test
	@echo "$(COLOR_GREEN)✓ All checks passed$(COLOR_RESET)"

## run: Build and run the application
run: build
	@echo "$(COLOR_BLUE)Running $(BINARY_NAME)...$(COLOR_RESET)"
	@./$(BINARY_NAME)

## clean: Remove binary and build artifacts
clean:
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE)
	@rm -f build.log
	@$(GO) clean -cache
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

## clean-data: Remove character save data
clean-data:
	@echo "$(COLOR_YELLOW)Removing character data...$(COLOR_RESET)"
	@rm -rf ~/.local/share/lazydndplayer/
	@echo "$(COLOR_GREEN)✓ Character data removed$(COLOR_RESET)"

## clean-run: Clean character data and run fresh
clean-run: clean-data build run

## clean-all: Remove everything (binary + data + coverage)
clean-all: clean clean-data
	@echo "$(COLOR_GREEN)✓ Everything cleaned$(COLOR_RESET)"

## install: Install the binary to $GOPATH/bin
install:
	@echo "$(COLOR_BLUE)Installing $(BINARY_NAME) to GOPATH/bin...$(COLOR_RESET)"
	@$(GO) install .
	@echo "$(COLOR_GREEN)✓ Installed successfully$(COLOR_RESET)"

# Docker targets (if needed in the future)
## docker-test: Run tests inside Docker
docker-test:
	@echo "$(COLOR_BLUE)Running tests in Docker...$(COLOR_RESET)"
	@if docker info > /dev/null 2>&1; then \
		echo "Docker is running"; \
		# Add docker test command here \
	else \
		echo "$(COLOR_RED)✗ Docker is not running$(COLOR_RESET)"; \
		echo "  Please start Docker Desktop and try again"; \
		exit 1; \
	fi

# Info targets
## info: Show project information
info:
	@echo "$(COLOR_BOLD)Project Information$(COLOR_RESET)"
	@echo "Binary name: $(BINARY_NAME)"
	@echo "Go version:  $$($(GO) version)"
	@echo "Build dir:   $$(pwd)"
	@echo ""
	@echo "$(COLOR_BOLD)Project Structure:$(COLOR_RESET)"
	@echo "  Source:    internal/"
	@echo "  Tests:     tests/"
	@echo "  Data:      data/"
	@echo ""
	@if [ -f $(BINARY_NAME) ]; then \
		echo "$(COLOR_BOLD)Binary Status:$(COLOR_RESET)"; \
		ls -lh $(BINARY_NAME); \
	else \
		echo "$(COLOR_YELLOW)Binary not built yet. Run: make build$(COLOR_RESET)"; \
	fi

## watch: Watch for changes and rebuild (requires entr)
watch:
	@echo "$(COLOR_BLUE)Watching for changes...$(COLOR_RESET)"
	@if command -v entr > /dev/null 2>&1; then \
		find . -name '*.go' | entr -c make build; \
	else \
		echo "$(COLOR_YELLOW)⚠ entr not installed$(COLOR_RESET)"; \
		echo "  Install: brew install entr (macOS) or apt install entr (Linux)"; \
	fi
