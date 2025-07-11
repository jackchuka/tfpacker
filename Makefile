.PHONY: build test clean install run

# Build variables
BINARY_NAME=tfpacker
BUILD_DIR=build

# Version information
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -X 'github.com/jackchuka/tfpacker/internal/version.Version=$(VERSION)'
LDFLAGS += -X 'github.com/jackchuka/tfpacker/internal/version.Commit=$(COMMIT)'
LDFLAGS += -X 'github.com/jackchuka/tfpacker/internal/version.Date=$(DATE)'

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

# Install the application
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install -ldflags "$(LDFLAGS)" .

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run golden tests
golden-test:
	@echo "Running golden tests..."
	@for dir in $$(find test -mindepth 1 -maxdepth 1 -type d); do \
		echo "Testing $${dir}..."; \
		rm -rf "$${dir}/golden"; \
		mkdir -p "$${dir}/golden"; \
		cd "$${dir}"; \
		tfpacker --config tfpacker.config.yaml --output golden; \
		cd ../..; \
		echo "✅ $${dir} golden generated"; \
	done

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run ./cmd/tfpacker $(ARGS)

# Run with dry-run mode
dry-run:
	@echo "Running in dry-run mode..."
	@go run ./cmd/tfpacker --dry-run $(ARGS)

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  install       - Install the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  golden-test   - Run golden tests for all test directories"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Run the application (use ARGS=\"arg1 arg2\" for arguments)"
	@echo "  dry-run       - Run in dry-run mode (use ARGS=\"arg1 arg2\" for arguments)"
	@echo "  help          - Show this help message"
