# Go-PWR Makefile

.PHONY: build clean install test fmt vet dev

# Build variables
BINARY_NAME=go-pwr
BUILD_DIR=build
CMD_DIR=cmd/go-pwr

# Version information
# For Windows compatibility, use simpler commands that work with Windows make
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>nul || echo unknown)
BUILD_DATE := $(shell powershell -NoProfile -Command "[DateTime]::UtcNow.ToString('yyyy-MM-ddTHH:mm:ssZ')" 2>nul || echo unknown)
LDFLAGS := -X main.gitCommit=$(GIT_COMMIT) -X main.buildDate=$(BUILD_DATE)

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"
	@mkdir -p $(BUILD_DIR)

	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)

	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	@GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)

	@echo "Cross-platform build complete!"

# Install the application
install:
	@echo "Installing $(BINARY_NAME)..."
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"
	@go install -ldflags "$(LDFLAGS)" ./$(CMD_DIR)
	@echo "Installation complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run the application in development mode
dev: build
	@echo "Running $(BINARY_NAME) in development mode..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Update dependencies
deps:
	@echo "Updating dependencies..."
	@go mod tidy
	@go mod download

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  install    - Install the application to GOPATH/bin"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  dev        - Build and run in development mode"
	@echo "  deps       - Update dependencies"
	@echo "  help       - Show this help message"
