# Project variables
PROJECT_NAME := packer-plugin-xelon

# Build variables
.DEFAULT_GOAL = test
BUILD_DIR := build
TOOLS_DIR := $(shell pwd)/tools
TOOLS_BIN_DIR := ${TOOLS_DIR}/bin
DEV_GOARCH := $(shell go env GOARCH)
DEV_GOOS := $(shell go env GOOS)
EXE =
ifeq ($(DEV_GOOS),windows)
EXE = .exe
endif
PATH := $(PATH):$(TOOLS_BIN_DIR)
SHELL := env PATH=$(PATH) /bin/bash


## tools: Install required tooling.
.PHONY: tools
tools:
	@echo "==> Installing required tooling..."
	@cd ${TOOLS_DIR} && GOBIN=${TOOLS_BIN_DIR} go install github.com/git-chglog/git-chglog/cmd/git-chglog
	@cd ${TOOLS_DIR} && GOBIN=${TOOLS_BIN_DIR} go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
	@cd ${TOOLS_DIR} && GOBIN=${TOOLS_BIN_DIR} go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc

## clean: Delete the build directory.
.PHONY: clean
clean:
	@echo "==> Removing '$(BUILD_DIR)' directory..."
	@rm -rf $(BUILD_DIR)

## lint: Lint code with golangci-lint.
.PHONY: lint
lint:
	@echo "==> Linting code with 'golangci-lint'..."
	@${TOOLS_BIN_DIR}/golangci-lint run

## test: Run all unit tests.
.PHONY: test
test:
	@echo "==> Running unit tests..."
	@mkdir -p $(BUILD_DIR)
	@go test -race -count=1 -v -cover -coverprofile=$(BUILD_DIR)/coverage.out -parallel=4 ./...

## testacc: Run all acceptance tests.
.PHONY: testacc
testacc:
	@echo "==> Running all acceptance tests..."
	@mkdir -p $(BUILD_DIR)
	@PACKER_ACC=1 go test -count=1 -v -cover -coverprofile=$(BUILD_DIR)/coverage-with-acceptance.out -parallel=4 -timeout 120m ./...

## generate: Generate necessary code and documentation for the plugin with packer-sdc.
.PHONY: generate
generate:
	@echo "==> Generating plugin code and documentation..."
	@go generate ./...
	@rm -rf .docs
	@${TOOLS_BIN_DIR}/packer-sdc renderdocs -src docs -partials docs-partials/ -dst .docs/
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs"
	@rm -r ".docs"

## build: Build plugin for default local system's operating system and architecture.
.PHONY: build
build:
	@echo "==> Building plugin..."
	@echo "    running go build for GOOS=$(DEV_GOOS) GOARCH=$(DEV_GOARCH)"
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME)$(EXE) main.go

## check: Checks plugin binary for compatibility with Packer.
.PHONY: check
check: build
	@echo "==> Checking plugin binary..."
	@# hack because plugin-check requires binary in the current directory
	@cp $(BUILD_DIR)/$(PROJECT_NAME)$(EXE) $(PROJECT_NAME)$(EXE)
	@${TOOLS_BIN_DIR}/packer-sdc plugin-check $(PROJECT_NAME)$(EXE)
	@rm $(PROJECT_NAME)$(EXE)


help: Makefile
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
