# Binary names
BINARY_NAME_LINUX=mr-valentine
BINARY_NAME_WINDOWS=mr-valentine.exe

# Build directory
BUILD_DIR=bin

.PHONY: all build build-linux build-windows clean

all: clean build

build: build-linux build-windows

build-linux:
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(BINARY_NAME_LINUX) ./cmd/app

build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(BINARY_NAME_WINDOWS) ./cmd/app

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)/*
	@go clean

run:
	@go run ./cmd/app 
