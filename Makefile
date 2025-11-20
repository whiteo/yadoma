# Yadoma Monorepo Makefile
.PHONY: test-agent test-webapp test-ui clean install-deps-agent install-deps-webapp install-deps-ui tools docker-build-agent docker-push-agent clean-docker build build-agent compress generate coverage-agent test-agent

# Variables
PROJECT_NAME=yadoma
GO_VERSION=1.25.1
NODE_VERSION=18
JAVA_VERSION=25

# Docker image/version
IMAGE_NAME ?= whiteo/yadoma
# Prefer explicit version (CI should pass YADOMA_AGENT_VERSION). Fallback to git describe or v0.1.0
VERSION ?= $(or $(YADOMA_AGENT_VERSION),$(shell git describe --tags --always --dirty 2>NUL),v0.1.0)

# Protobuf tool versions (pin exact versions for reproducibility)
PROTOC_GEN_GO_VERSION ?= v1.34.2
PROTOC_GEN_GO_GRPC_VERSION ?= v1.5.1

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

# =============================================================================
# DEPENDENCY MANAGEMENT
# =============================================================================

install-deps: ## Install Go dependencies for agent
	@echo "$(BLUE)📦 Installing Go dependencies for agent...$(NC)"
	cd agent && go mod download && go mod verify
	@echo "$(GREEN)✅ Go dependencies installed$(NC)"
# =============================================================================
# TESTING
# =============================================================================

test-agent: ## Run Go agent tests
	@echo "$(BLUE)🧪 Running Go agent tests...$(NC)"
	cd agent && go test -v -race -count=1 ./...
	@echo "$(GREEN)✅ Go agent tests passed$(NC)"

coverage-agent: ## Run Go agent tests with coverage
	@echo "$(BLUE)🧪 Running Go agent tests with coverage...$(NC)"
	cd agent && go test -v -race -count=1 -coverprofile=coverage.txt ./...
	@echo "$(GREEN)✅ Coverage report generated at agent/coverage.txt$(NC)"

# =============================================================================
# BUILD
# =============================================================================

build-agent: ## Build Go agent binary
	@echo "$(BLUE)🔨 Building Go agent...$(NC)"
	mkdir -p build
	cd agent && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../build/yadoma-agent ./cmd
	@echo "$(GREEN)✅ Go agent built successfully$(NC)"

compress: ## Compress binary with UPX
	@echo "$(BLUE)📦 Compressing binary with UPX...$(NC)"
	upx --best --lzma ./build/yadoma-agent
	@echo "$(GREEN)✅ Binary compressed$(NC)"

build-ui: ## Build TypeScript/React UI
	@echo "$(BLUE)🎨 Building UI (React/Vite)...$(NC)"
	cd ui && npm install && npm run build
	@echo "$(GREEN)✅ UI built successfully$(NC)"

build-webapp: ## Build Java webapp with embedded UI
	@echo "$(BLUE)☕ Building Java webapp...$(NC)"
	cd webapp && chmod +x gradlew && ./gradlew clean build -x test
	@echo "$(GREEN)✅ Webapp built successfully$(NC)"

test-webapp: ## Run webapp tests with coverage
	@echo "$(BLUE)🧪 Running webapp tests...$(NC)"
	cd webapp && chmod +x gradlew && ./gradlew test jacocoTestReport
	@echo "$(GREEN)✅ Webapp tests passed$(NC)"

build: build-agent compress ## Build all components (currently only Go agent)
	@echo "$(GREEN)✅ All builds completed$(NC)"

# =============================================================================
# DOCKER
# =============================================================================

docker-build-agent: ## Build Docker image for Go agent with explicit tag
	@echo "$(BLUE)🐳 Building Docker image $(IMAGE_NAME):$(VERSION)...$(NC)"
	docker build -f Dockerfile -t $(IMAGE_NAME):$(VERSION) --build-arg DOCKER_GID=$(DOCKER_GID) .
	@echo "$(GREEN)✅ Docker image built: $(IMAGE_NAME):$(VERSION)$(NC)"

docker-push-agent: docker-build-agent ## Push Docker image with explicit tag
	@echo "$(BLUE)📤 Pushing Docker image $(IMAGE_NAME):$(VERSION)...$(NC)"
	docker push $(IMAGE_NAME):$(VERSION)
	@echo "$(GREEN)✅ Docker image pushed: $(IMAGE_NAME):$(VERSION)$(NC)"

# =============================================================================
# CLEANING
# =============================================================================

clean: ## Clean build artifacts and coverage reports
	@echo "$(BLUE)🧹 Cleaning build artifacts...$(NC)"
	rm -rf build/
	rm -rf agent/coverage.out agent/coverage.html agent/coverage.txt
	rm -rf webapp/build webapp/src/main/resources/static/*
	rm -rf ui/dist
	@echo "$(GREEN)✅ Clean completed$(NC)"

clean-docker: ## Clean Docker images (by explicit tag)
	@echo "$(BLUE)🐳 Cleaning Docker image $(IMAGE_NAME):$(VERSION)...$(NC)"
	-docker rmi $(IMAGE_NAME):$(VERSION)
	@echo "$(GREEN)✅ Docker cleanup completed$(NC)"

# =============================================================================
# ALL-IN-ONE COMMANDS
# =============================================================================
all: clean install-deps tools test-agent coverage-agent generate build compress ## Full pipeline
	@echo "$(GREEN)🚀 All tasks completed successfully!$(NC)"