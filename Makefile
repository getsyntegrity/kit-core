.PHONY: all deps deps-update build test

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
