# UFO File Upload Web Application Makefile

# Variables
BINARY_NAME=ufo-server
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server/main.go
DATA_DIR=./data

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building UFO server..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

# Run the application
.PHONY: run
run: build
	@echo "Starting UFO server..."
	@mkdir -p $(DATA_DIR)
	@$(BINARY_PATH)

# Run in development mode with auto-reload (requires air)
.PHONY: dev
dev:
	@if command -v air > /dev/null; then \
		echo "Starting UFO in development mode..."; \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running normally..."; \
		make run; \
	fi

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Build for production (optimized)
.PHONY: build-prod
build-prod:
	@echo "Building for production..."
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Production build complete: $(BINARY_PATH)"

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@echo "Development tools installed"

# Create data directory
.PHONY: setup
setup:
	@echo "Setting up UFO..."
	@mkdir -p $(DATA_DIR)
	@mkdir -p web/static/css
	@mkdir -p web/static/js
	@mkdir -p web/static/images
	@mkdir -p web/templates
	@echo "Directory structure created"

# Build and run with custom port
.PHONY: run-port
run-port: build
	@echo "Starting UFO server on port $(PORT)..."
	@mkdir -p $(DATA_DIR)
	@PORT=$(PORT) $(BINARY_PATH)

# Show help
.PHONY: help
help:
	@echo "UFO File Upload Server - Available Commands:"
	@echo ""
	@echo "  make build       - Build the application"
	@echo "  make run         - Build and run the server"
	@echo "  make dev         - Run in development mode (requires air)"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make deps        - Download and tidy dependencies"
	@echo "  make test        - Run tests"
	@echo "  make build-prod  - Build optimized production binary"
	@echo "  make setup       - Create directory structure"
	@echo "  make install-tools - Install development tools"
	@echo "  make run-port PORT=8081 - Run on custom port"
	@echo "  make help        - Show this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  PORT         - Server port (default: 8080)"
	@echo "  DATA_DIR     - Upload directory (default: ./data)"
	@echo "  MAX_FILE_SIZE - Max file size in bytes (default: 104857600 = 100MB)"