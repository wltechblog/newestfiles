# Makefile for newestfiles

# Variables
BINARY_NAME=newestfiles
GO_FILES=$(wildcard *.go)
BUILD_DIR=build

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(GO_FILES)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet
.PHONY: vet
vet:
	go vet ./...

# Run all checks (fmt, vet, test)
.PHONY: check
check: fmt vet test

# Install the binary to GOPATH/bin
.PHONY: install
install:
	go install .

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the program with example arguments
.PHONY: run-example
run-example: build
	./$(BUILD_DIR)/$(BINARY_NAME) .go .txt

# Run the program with JSON output
.PHONY: run-json
run-json: build
	./$(BUILD_DIR)/$(BINARY_NAME) -j .go .txt

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format Go code"
	@echo "  vet           - Run go vet"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  clean         - Clean build artifacts"
	@echo "  run-example   - Run with example arguments"
	@echo "  run-json      - Run with JSON output"
	@echo "  help          - Show this help message"
