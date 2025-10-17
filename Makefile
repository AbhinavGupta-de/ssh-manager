# SSH Manager Makefile

APP_NAME=sshm
VERSION=1.0.0
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Build the binary for current platform
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	$(GOBUILD) -o $(APP_NAME) -v

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64
	
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64
	
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64
	
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64
	
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe
	
	@echo "✅ All builds complete! Check the $(BUILD_DIR) directory"

# Install locally (Unix-like systems)
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	@chmod +x $(APP_NAME)
	@if [ -w /usr/local/bin ]; then \
		mv $(APP_NAME) /usr/local/bin/; \
		echo "✅ Installed to /usr/local/bin/$(APP_NAME)"; \
	elif [ -d $$HOME/.local/bin ]; then \
		mv $(APP_NAME) $$HOME/.local/bin/; \
		echo "✅ Installed to $$HOME/.local/bin/$(APP_NAME)"; \
	else \
		mkdir -p $$HOME/bin; \
		mv $(APP_NAME) $$HOME/bin/; \
		echo "✅ Installed to $$HOME/bin/$(APP_NAME)"; \
	fi

# Uninstall
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	@rm -f /usr/local/bin/$(APP_NAME)
	@rm -f $$HOME/.local/bin/$(APP_NAME)
	@rm -f $$HOME/bin/$(APP_NAME)
	@echo "✅ Uninstalled"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(APP_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run tests (when you add them)
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run the application
.PHONY: run
run:
	@$(GOBUILD) -o $(APP_NAME) -v
	@./$(APP_NAME)

# Display help
.PHONY: help
help:
	@echo "SSH Manager - Build Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build for current platform"
	@echo "  make build-all      Build for all platforms"
	@echo "  make install        Build and install locally"
	@echo "  make uninstall      Remove installed binary"
	@echo "  make clean          Clean build artifacts"
	@echo "  make test           Run tests"
	@echo "  make run            Build and run"
	@echo "  make help           Display this help message"

# Default target
.DEFAULT_GOAL := help