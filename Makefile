.PHONY: build install test clean

# Get version from VERSION file
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")

# Build binary
build:
	@echo "🔨 Building AI Terminal Helper v$(VERSION)..."
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/ai-helper ./cmd/ai-helper
	@echo "✅ Build complete: bin/ai-helper"

# Install to ~/.ai/ (clean install)
install: build
	@echo "📦 Installing AI Terminal Helper..."
	@echo "🧹 Cleaning old files..."
	@rm -f ~/.ai/ai-helper.sh ~/.ai/cache-manager.sh ~/.ai/zsh-integration.sh
	@mkdir -p ~/.ai
	@cp bin/ai-helper ~/.ai/ai-helper
	@chmod +x ~/.ai/ai-helper
	@cp integrations/zsh/ai-helper.zsh ~/.ai/
	@echo "✅ Installed to ~/.ai/"
	@echo ""
	@echo "📝 Next steps:"
	@echo "  1. Add to ~/.zshrc: echo 'source ~/.ai/ai-helper.zsh' >> ~/.zshrc"
	@echo "  2. Reload shell: source ~/.zshrc"
	@echo "  3. Test: ask how do I list all pods"
	@echo ""
	@echo "💡 Note: ~/.ai is automatically added to PATH by ai-helper.zsh"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Clean build artifacts and installation
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@echo "✅ Clean complete"

# Uninstall from ~/.ai/
uninstall:
	@echo "🗑️  Uninstalling AI Terminal Helper..."
	@rm -f ~/.ai/ai-helper
	@rm -f ~/.ai/ai-helper.zsh
	@rm -f ~/.ai/ai-helper.sh
	@rm -f ~/.ai/cache-manager.sh
	@rm -f ~/.ai/zsh-integration.sh
	@echo "✅ Uninstalled from ~/.ai/"
	@echo ""
	@echo "📝 Don't forget to remove 'source ~/.ai/ai-helper.zsh' from ~/.zshrc"

# Build for multiple platforms
build-all:
	@echo "🔨 Building for all platforms v$(VERSION)..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/ai-helper-darwin-amd64 ./cmd/ai-helper
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/ai-helper-darwin-arm64 ./cmd/ai-helper
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/ai-helper-linux-amd64 ./cmd/ai-helper
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/ai-helper-linux-arm64 ./cmd/ai-helper
	@echo "✅ Build complete for all platforms"

# Check dependencies
deps:
	@echo "📦 Checking dependencies..."
	@go mod tidy
	@go mod verify
	@echo "✅ Dependencies OK"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run ./...

# Show version
version:
	@./bin/ai-helper version || echo "Build first: make build"

