# Project
BINARY_NAME := {{.Inputs.project_name.Value}}
MODULE := github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}
{{- if eq .Inputs.project_type.Value "cli"}}

# Release targets
RELEASE_OS := linux darwin
RELEASE_ARCH := amd64 arm64
{{- end}}

.PHONY: all clean build test lint fmt tidy deps{{if eq .Inputs.project_type.Value "cli"}} release $(RELEASE_OS){{else}} run docker{{end}}

# Build directories
BUILD_DIR := ${PWD}/build
BINARY_DIR := ${BUILD_DIR}/bin

# Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# Build flags
export CGO_ENABLED := 0
GCFLAGS := -gcflags "all=-trimpath=${PWD}"
ASMFLAGS := -asmflags "all=-trimpath=${PWD}"
LDFLAGS := -ldflags "-X ${MODULE}/internal/version.Version=${VERSION} -s -w"

GOLANGCI_LINT_VERSION := v2.5.0

# Default target
all: lint test build

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ${BUILD_DIR}
	@go clean

${BUILD_DIR}:
	@mkdir -p ${BINARY_DIR}

# Run tests with coverage (the race detector needs cgo, so it is switched on only here)
test: ${BUILD_DIR}
	@echo "Running tests..."
	@CGO_ENABLED=1 go test -v -race -coverprofile=${BUILD_DIR}/coverage.out ./...
	@go tool cover -html=${BUILD_DIR}/coverage.out -o ${BUILD_DIR}/coverage.html

# Build for current platform
build: ${BUILD_DIR}
	@echo "Building binary..."
	@go build ${GCFLAGS} ${ASMFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/${BINARY_NAME}
{{- if eq .Inputs.project_type.Value "cli"}}

# Build for all platforms
release: ${BUILD_DIR} test $(RELEASE_OS)
	@echo "Built release artifacts in ${BINARY_DIR}"

$(RELEASE_OS):
	@for arch in $(RELEASE_ARCH); do \
		echo "Building for $@/$$arch..."; \
		GOOS=$@ GOARCH=$$arch go build ${GCFLAGS} ${ASMFLAGS} ${LDFLAGS} \
			-o ${BINARY_DIR}/${BINARY_NAME}_$@_$$arch ./cmd/${BINARY_NAME}; \
	done
{{- else}}

# Run the service locally
run:
	@go run ${LDFLAGS} ./cmd/${BINARY_NAME}

# Build the container image
docker:
	@docker build --build-arg VERSION=${VERSION} -t ${BINARY_NAME}:${VERSION} .
{{- end}}

# Run linting
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Please run: make deps"; \
		exit 1; \
	fi

# Format code
fmt:
	@golangci-lint fmt

# Tidy go.mod
tidy:
	@go mod tidy

# Install dependencies and tooling
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
