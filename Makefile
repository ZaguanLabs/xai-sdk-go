SHELL := /bin/bash

PROJECT := github.com/ZaguanLabs/xai-sdk-go

# Find all main packages in the examples directory
EXAMPLES := $(shell find examples -type f -name "*.go" -exec grep -l "package main" {} +)

.PHONY: help fmt lint test test-integration proto clean examples

help:
	@echo "Common targets"
	@echo "  fmt              - gofmt over the module"
	@echo "  lint             - placeholder lint target"
	@echo "  test             - run unit tests"
	@echo "  test-integration - run integration tests (requires XAI_API_KEY)"
	@echo "  proto            - regenerate protobuf bindings"
	@echo "  examples         - build and list all examples"
	@echo "  clean            - remove build artifacts"

fmt:
	@gofmt -w $(shell find . -name '*.go')

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

test:
	@go test ./...

test-integration:
	@echo "Running integration tests (requires XAI_API_KEY)..."
	@if [ -z "$$XAI_API_KEY" ]; then \
		echo "Error: XAI_API_KEY environment variable not set"; \
		echo "Set it with: export XAI_API_KEY=your-api-key"; \
		exit 1; \
	fi
	@go test -tags=integration -v ./xai/embed ./xai/files ./xai/image ./xai/auth

proto:
	@$(HOME)/go/bin/buf generate proto

examples:
	@echo "Building and listing examples..."
	@for ex in $(EXAMPLES); do \
		echo "  - $$ex"; \
		go build -o /dev/null $$ex; \
	done

clean:
	@rm -rf bin build dist
