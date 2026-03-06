## Constants
# Linter version
LINT_VERSION := v2.10.1


# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

unit_test: ## Run all unit tests in the pkg folder
	@echo "Running unit test framework..."
	go test -v ./pkg/... ./internal/...

unit_test_cover: ## Run all unit tests in the pkg folder
	@echo "Running unit test framework with coverage..."
	go test -cover ./pkg/... ./internal/...

lint_local_clean: ## Cleanup the local cache from the linter
	@echo "Cleaning linter cache..."
	golangci-lint cache clean

lint_local: lint_local_clean ## Run local instance of linting across the code base
	@echo "Running linter on codebase..."
	golangci-lint run ./pkg/... ./internal/...

lint_local_fix: ## Run local instance of linting across the code base including auto-fixing
	@echo "Running linter with fix option..."
	golangci-lint run --fix ./pkg/... ./internal/...

lint_docker: ## Run docker instance of linting across the code base
	docker run --rm -t -v $(PWD):/app -v ~/.cache/golangci-lint/$(LINT_VERSION):/root/.cache -w /app golangci/golangci-lint:$(LINT_VERSION) golangci-lint run -v ./pkg/... ./internal/...

lint_docker_fix: ## Run docker instance of linting across the code base including auto-fixing
	docker run --rm -v $(PWD):/app -v ~/.cache/golangci-lint/$(LINT_VERSION):/root/.cache -w /app golangci/golangci-lint:$(LINT_VERSION) golangci-lint run --fix ./pkg/... ./internal/...
