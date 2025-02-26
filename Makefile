# Binary names
BINARY_NAME_LINUX=mr-valentine
BINARY_NAME_WINDOWS=mr-valentine.exe

# Build directory
BUILD_DIR=bin

REQUIRED_EXECUTABLES := go npm
$(foreach exec,$(REQUIRED_EXECUTABLES),\
    $(if $(shell which $(exec)),,$(error "$(exec) not found in PATH. Please install $(exec): $(if $(filter go,$(exec)),https://golang.org/doc/install,https://nodejs.org/)")))

.PHONY: all build build-linux build-windows clean frontend run dev install-tools

all: clean frontend build

build: build-linux build-windows

format:
	@echo "Formatting..."
	@go fmt ./...

build-linux:
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(BINARY_NAME_LINUX) ./cmd/app

build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(BINARY_NAME_WINDOWS) ./cmd/app

frontend:
	@echo "Building frontend assets..."
	@npm run build

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)/*
	@rm -rf ui/static/scripts/dist/*
	@rm -rf tmp/*
	@go clean

run: frontend
	@go run ./cmd/app 

install-tools:
	@echo "Checking prerequisites..."
	@which go >/dev/null || (echo "Error: Go is not installed. Please install Go first: https://golang.org/doc/install" && exit 1)
	@which npm >/dev/null || (echo "Error: npm is not installed. Please install Node.js and npm first: https://nodejs.org/" && exit 1)
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@npm install

dev: install-tools
	@echo "Starting development mode..."
	@npm run watch & air
