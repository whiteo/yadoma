# Yadoma Monorepo Makefile
.PHONY: test-agent test-webapp test-ui clean install-deps-agent install-deps-webapp install-deps-ui

# Variables
PROJECT_NAME=yadoma
GO_VERSION=1.25.1
NODE_VERSION=18
JAVA_VERSION=25

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

# =============================================================================
# DEPENDENCY MANAGEMENT
# =============================================================================

install-deps-agent: ## Install Go dependencies for agent
	@echo "$(BLUE)📦 Installing Go dependencies for agent...$(NC)"
	cd agent && go mod download && go mod verify
	@echo "$(GREEN)✅ Go dependencies installed$(NC)"

install-deps-webapp: ## Install Java dependencies for webapp (placeholder)
	@echo "$(YELLOW)⚠️  Webapp (Java) dependencies - not implemented yet$(NC)"
	# cd webapp && ./gradlew dependencies || cd webapp && mvn dependency:resolve

install-deps-ui: ## Install Node.js dependencies for UI (placeholder)
	@echo "$(YELLOW)⚠️  UI (TypeScript/React) dependencies - not implemented yet$(NC)"
	# cd ui && npm install

tools: ## Install dev tools
	@echo "$(BLUE)🔨 Installing dev tools...$(NC)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "$(GREEN)✅ Tools installed$(NC)"

install-deps: install-deps-agent tools ## Install all dependencies
		@echo "$(GREEN)✅ All dependencies installed$(NC)"

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

test-webapp: ## Run Java webapp tests (placeholder)
	@echo "$(YELLOW)⚠️  Webapp (Java) tests - not implemented yet$(NC)"
	# cd webapp && ./gradlew test || cd webapp && mvn test

test-ui: ## Run TypeScript/React UI tests (placeholder)
	@echo "$(YELLOW)⚠️  UI (TypeScript/React) tests - not implemented yet$(NC)"
	# cd ui && npm test

# =============================================================================
# CODE GENERATION
# =============================================================================

generate: ## Generate protobuf files for Go agent
	@echo "$(BLUE)🔧 Generating protobuf files...$(NC)"
	go generate ./generate.go
	@echo "$(GREEN)✅ Protobuf files generated$(NC)"

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

build-webapp: ## Build Java webapp (placeholder)
	@echo "$(YELLOW)⚠️  Webapp (Java) build - not implemented yet$(NC)"
	# cd webapp && ./gradlew build || cd webapp && mvn package

build-ui: ## Build TypeScript/React UI (placeholder)
	@echo "$(YELLOW)⚠️  UI (TypeScript/React) build - not implemented yet$(NC)"
	# cd ui && npm run build

build: build-agent compress ## Build all components (currently only Go agent)
	@echo "$(GREEN)✅ All builds completed$(NC)"

# =============================================================================
# DOCKER
# =============================================================================

docker-build-agent: build-agent compress ## Build Docker image for Go agent
	@echo "$(BLUE)🐳 Building Docker image for Go agent...$(NC)"
	docker build -f Dockerfile -t $(PROJECT_NAME)/agent:latest .
	@echo "$(GREEN)✅ Docker image built$(NC)"

docker-push-agent: docker-build-agent ## Push Docker image to Docker Hub
	@echo "$(BLUE)📤 Pushing Docker image to Docker Hub...$(NC)"
	docker tag $(PROJECT_NAME)/agent:latest $(DOCKER_USER)/$(PROJECT_NAME)-agent:latest
	docker tag $(PROJECT_NAME)/agent:latest $(DOCKER_USER)/$(PROJECT_NAME)-agent:$(GITHUB_SHA)
	docker push $(DOCKER_USER)/$(PROJECT_NAME)-agent:latest
	docker push $(DOCKER_USER)/$(PROJECT_NAME)-agent:$(GITHUB_SHA)
	@echo "$(GREEN)✅ Docker image pushed$(NC)"

# =============================================================================
# CLEANING
# =============================================================================

clean: ## Clean build artifacts and coverage reports
	@echo "$(BLUE)🧹 Cleaning build artifacts...$(NC)"
	rm -rf build/
	rm -rf agent/coverage.out agent/coverage.html
	# rm -rf webapp/build webapp/target
	# rm -rf ui/build ui/dist ui/node_modules/.cache
	@echo "$(GREEN)✅ Clean completed$(NC)"

clean-docker: ## Clean Docker containers and images
	@echo "$(BLUE)🐳 Cleaning Docker resources...$(NC)"
	docker stop $(PROJECT_NAME)-agent || true
	docker rm $(PROJECT_NAME)-agent || true
	docker rmi $(PROJECT_NAME)/agent:latest || true
	@echo "$(GREEN)✅ Docker cleanup completed$(NC)"
