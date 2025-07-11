.PHONY: build test clean install run

# Build variables
BINARY_NAME=tfpacker
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=build
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Install the application
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) .

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
