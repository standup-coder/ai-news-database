# News4Coder Makefile
# Quick reference for development tasks

BINARY_NAME=news4coder
MAIN_PATH=main.go
BUILD_DIR=.

# Version info (override via make build VERSION=v1.0.0)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -ldflags "-X github.com/spf13/cobra.version=$(VERSION) -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

.PHONY: all build run test clean install fmt vet lint release help

all: build

## build: Build the binary for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

## run: Build and run the binary
run: build
	./$(BINARY_NAME)

## test: Run all tests with coverage
test:
	go test ./... -v -race -cover

## test-short: Run tests quickly
 test-short:
	go test ./... -short

## clean: Remove built binaries and test artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	go clean -testcache

## install: Install binary to $GOPATH/bin
install:
	go install $(LDFLAGS) .

## fmt: Format Go code
fmt:
	go fmt ./...

## vet: Run go vet
vet:
	go vet ./...

## lint: Run golangci-lint (requires installation)
lint:
	golangci-lint run ./...

## mod: Download and tidy dependencies
mod:
	go mod download
	go mod tidy

## release: Build cross-platform binaries for release
release:
	@echo "Building release binaries..."
	mkdir -p dist
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Done. Binaries in dist/"

## web: Serve the web UI locally on port 8080
web:
	cd web && python3 -m http.server 8080

## docs: Check all markdown files are present
docs:
	@echo "Checking documentation..."
	@ls -1 *.md

## help: Show this help message
help:
	@cat Makefile | grep -E "^##" | sed 's/## //'
