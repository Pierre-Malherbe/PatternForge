.PHONY: build run install clean test help

BINARY_NAME=patternforge
BUILD_DIR=bin
INSTALL_PATH=/usr/local/bin

# Build for current platform
build:
	@echo "🔨 Building PatternForge..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/patternforge
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for macOS (Apple Silicon and Intel)
build-macos:
	@echo "🍎 Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/patternforge
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/patternforge
	@echo "✅ macOS builds complete"

# Run without building
run:
	@go run ./cmd/patternforge

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Install to system
install: build
	@echo "📥 Installing to $(INSTALL_PATH)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@echo "✅ Installed! Run with: patternforge"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Show help
help:
	@echo "PatternForge - Build Commands"
	@echo ""
	@echo "  make build        Build for current platform"
	@echo "  make build-macos  Build for macOS (ARM64 + AMD64)"
	@echo "  make run          Run without building"
	@echo "  make deps         Install Go dependencies"
	@echo "  make install      Install to /usr/local/bin"
	@echo "  make clean        Remove build artifacts"
	@echo "  make test         Run tests"
	@echo "  make help         Show this help"
