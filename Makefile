.PHONY: all build test format lint clean run help

# Default target
all: format lint test build

# Build the application
build:
	go build -o claude-usage-go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Format code
format:
	@echo "Running gofmt..."
	@gofmt -w .
	@echo "Running goimports (if available)..."
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "goimports not installed, skipping"

# Run linters
lint:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Running golint (if available)..."
	@command -v golint >/dev/null 2>&1 && golint ./... || echo "golint not installed, skipping"

# Clean build artifacts
clean:
	rm -f claude-usage-go
	go clean

# Run the application
run: build
	./claude-usage-go daily

# Install development tools
install-tools:
	@echo "Installing goimports..."
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing golint..."
	go install golang.org/x/lint/golint@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Format, lint, test, and build (default)"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  format        - Format code with gofmt and goimports"
	@echo "  lint          - Run linters (go vet and golint)"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Build and run the application"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help message"