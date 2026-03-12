.PHONY: all deps deps-update build test lint

# Default: run deps, build and test
all: deps build test

# Install / refresh dependencies (download, tidy)
deps:
	go mod download
	go mod tidy

# Update all dependencies to latest minor/patch (go get -u ./...)
deps-update:
	go get -u ./...
	go mod tidy

# Build all packages
build:
	go build ./...

# Run tests
test:
	go test ./...

# Run golangci-lint (correct path for v2: main is in cmd/golangci-lint).
# Use this before forge pr until devforge fixes its static-analysis invocation.
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.0 run --timeout=5m
