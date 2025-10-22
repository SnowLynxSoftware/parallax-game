# Makefile for parallax-game
# Best practices Go project automation

# Variables
APP_NAME := parallax-game
DOCKER_IMAGE := ghcr.io/snowlynxsoftware/$(APP_NAME)
DOCKER_TAG := latest
COMPOSE_FILE := docker-compose.yml
GO_VERSION := 1.25
BUILD_DIR := build
COVERAGE_DIR := coverage

# Go related variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOLINT := golangci-lint

# Build flags
LDFLAGS := -ldflags="-w -s"
BUILD_FLAGS := -v $(LDFLAGS)

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

.PHONY: help
help: ## Display this help message
	@echo "$(CYAN)Available commands:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: deps
deps: ## Download and install dependencies
	@echo "$(CYAN)Downloading dependencies...$(RESET)"
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "$(CYAN)Updating dependencies...$(RESET)"
	$(GOGET) -u ./...
	$(GOMOD) tidy

.PHONY: build
build: clean deps ## Build the application
	@echo "$(CYAN)Building $(APP_NAME)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(APP_NAME)$(RESET)"

.PHONY: build-local
build-local: clean deps ## Build the application for local development
	@echo "$(CYAN)Building $(APP_NAME) for local development...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "$(GREEN)Local build completed: $(BUILD_DIR)/$(APP_NAME)$(RESET)"

.PHONY: run
run: ## Run the application locally
	@echo "$(CYAN)Running $(APP_NAME)...$(RESET)"
	$(GOCMD) run .

.PHONY: migrate
migrate: ## Run database migrations
	@echo "$(CYAN)Running database migrations...$(RESET)"
	$(GOCMD) run . migrate

.PHONY: test
test: ## Run all tests with summary
	@echo -e "$(CYAN)Running tests...$(RESET)"
	@bash -c 'set -o pipefail; \
	OUTPUT=$$(mktemp); \
	$(GOTEST) ./... -v 2>&1 | tee $$OUTPUT; \
	EXIT=$$?; \
	PASSED=$$(grep "^--- PASS:" $$OUTPUT | wc -l); \
	FAILED=$$(grep "^--- FAIL:" $$OUTPUT | wc -l); \
	TOTAL=$$((PASSED + FAILED)); \
	echo ""; \
	echo -e "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"; \
	echo -e "$(CYAN)Test Summary:$(RESET)"; \
	echo -e "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"; \
	echo -e "  $(GREEN)✓ Passed:$(RESET)  $$PASSED"; \
	if [ $$FAILED -gt 0 ]; then \
		echo -e "  $(RED)✗ Failed:$(RESET)  $$FAILED"; \
	else \
		echo -e "  $(GREEN)✗ Failed:$(RESET)  $$FAILED"; \
	fi; \
	echo -e "  $(CYAN)Total:$(RESET)    $$TOTAL"; \
	echo -e "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"; \
	rm -f $$OUTPUT; \
	if [ $$FAILED -gt 0 ]; then \
		echo ""; \
		echo -e "$(RED)Tests failed!$(RESET)"; \
		exit 1; \
	else \
		echo ""; \
		echo -e "$(GREEN)All tests passed!$(RESET)"; \
	fi; \
	exit $$EXIT'

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(CYAN)Running tests with coverage...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) ./... -v -coverprofile=$(COVERAGE_DIR)/coverage.out
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_DIR)/coverage.html$(RESET)"

.PHONY: benchmark
benchmark: ## Run benchmark tests
	@echo "$(CYAN)Running benchmark tests...$(RESET)"
	$(GOTEST) ./... -bench=. -benchmem


.PHONY: fmt
fmt: ## Format Go code
	@echo "$(CYAN)Formatting code...$(RESET)"
	$(GOFMT) -s -w .

.PHONY: vet
vet: ## Run go vet
	@echo "$(CYAN)Running go vet...$(RESET)"
	$(GOCMD) vet ./...

.PHONY: lint
lint: ## Run golangci-lint (if installed)
	@echo "$(CYAN)Running golangci-lint...$(RESET)"
	@which $(GOLINT) > /dev/null || (echo "$(YELLOW)golangci-lint not installed. Skipping. Install from https://golangci-lint.run/$(RESET)" && exit 0)
	$(GOLINT) run ./...

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)

.PHONY: docker-build
docker-build: ## Build Docker image for production
	@echo "$(CYAN)Building Docker image...$(RESET)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(RESET)"

.PHONY: docker-push
docker-push: docker-build ## Build and push Docker image to registry
	@echo "$(CYAN)Pushing Docker image...$(RESET)"
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)Docker image pushed: $(DOCKER_IMAGE):$(DOCKER_TAG)$(RESET)"

.PHONY: compose-up
compose-up: ## Start services with docker-compose in detached mode
	@echo "$(CYAN)Starting services with docker-compose...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo "$(GREEN)Services started in detached mode$(RESET)"

.PHONY: compose-down
compose-down: ## Stop and remove docker-compose services
	@echo "$(CYAN)Stopping docker-compose services...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) down
	@echo "$(GREEN)Services stopped$(RESET)"

.PHONY: compose-logs
compose-logs: ## Show docker-compose logs
	@echo "$(CYAN)Showing docker-compose logs...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) logs -f

.PHONY: compose-build
compose-build: ## Build docker-compose services
	@echo "$(CYAN)Building docker-compose services...$(RESET)"
	docker-compose -f $(COMPOSE_FILE) build

.PHONY: compose-restart
compose-restart: compose-down compose-up ## Restart docker-compose services

.PHONY: clean
clean: ## Clean build artifacts and cache
	@echo "$(CYAN)Cleaning...$(RESET)"
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@echo "$(GREEN)Cleaned$(RESET)"

.PHONY: clean-docker
clean-docker: ## Clean Docker containers and images
	@echo "$(CYAN)Cleaning Docker containers and images...$(RESET)"
	docker system prune -f
	@echo "$(GREEN)Docker cleanup completed$(RESET)"

.PHONY: status
status: ## Show project status
	@echo "$(CYAN)Project Status:$(RESET)"
	@echo "App Name: $(APP_NAME)"
	@echo "Go Version: $(GO_VERSION)"
	@echo ""
	@echo "$(CYAN)Git Status:$(RESET)"
	@git status --porcelain
	@echo ""
	@echo "$(CYAN)Docker Compose Status:$(RESET)"
	@docker-compose -f $(COMPOSE_FILE) ps

# Default target
.DEFAULT_GOAL := help
