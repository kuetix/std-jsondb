APP_NAME := jsondb
BUILD_DIR := runtime/bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'

.DEFAULT_GOAL := help

.PHONY: help all build build-pkg clean test tag

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the decision package runner
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-pkg ./cmd/pkg

all: build ## Build the package

build-pkg: ## Build the decision package runner
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-pkg ./cmd/pkg

tag: ## Create an annotated release tag (usage: make tag TAG=v0.x.y)
	@if [ -z "$(TAG)" ]; then echo "Usage: make tag TAG=v0.x.y"; exit 1; fi
	@case "$(TAG)" in v[0-9]*) ;; *) echo "TAG must look like vX.Y.Z"; exit 1;; esac
	@if [ -n "$$(git status --porcelain)" ]; then echo "working tree not clean"; exit 1; fi
	git tag -a $(TAG) -m "Release $(TAG)"
	@echo "Created tag $(TAG). Push it with:  git push origin $(TAG)"
	@echo "The release workflow builds and publishes the GitHub Release."

test: ## Run package tests
	go test ./...

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
