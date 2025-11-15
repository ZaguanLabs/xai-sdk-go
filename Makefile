SHELL := /bin/bash

PROJECT := github.com/ZaguanLabs/xai-sdk-go

# Find all main packages in the examples directory
EXAMPLES := $(shell find examples -type f -name "*.go" -exec grep -l "package main" {} +)

.PHONY: help fmt lint test proto clean examples

help:
	@echo "Common targets"
	@echo "  fmt       - gofmt over the module"
	@echo "  lint      - placeholder lint target"
	@echo "  test      - run unit tests"
	@echo "  proto     - regenerate protobuf bindings"
	@echo "  examples  - build and list all examples"
	@echo "  clean     - remove build artifacts"

fmt:
	@gofmt -w $(shell find . -name '*.go')

lint:
	@echo "lint: configure golangci-lint in later phase"

test:
	@go test ./...

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
