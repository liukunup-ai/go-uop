# go-uop Makefile

.PHONY: help
help: ## Show this help message
	@echo "go-uop Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-][a-zA-Z0-9/_-]*:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Build Commands
# =============================================================================
.PHONY: build
build: ## Build all packages
	go build ./...

.PHONY: build/ios
build/ios: ## Build iOS module only
	go build ./ios/...

.PHONY: build/android
build/android: ## Build Android module only
	go build ./android/...

.PHONY: build/yaml
build/yaml: ## Build YAML runner only
	go build ./yaml/...

# =============================================================================
# Test Commands
# =============================================================================
.PHONY: test
test: ## Run all tests
	go test -v ./...

.PHONY: test/short
test/short: ## Run tests with -short flag
	go test -v -short ./...

.PHONY: test/cover
test/cover: ## Run tests with coverage report
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test/locator
test/locator: ## Run locator tests only
	go test -v ./internal/locator/...

.PHONY: test/yaml
test/yaml: ## Run YAML parser tests only
	go test -v ./yaml/...

.PHONY: test/ios
test/ios: ## Run iOS tests only
	go test -v ./ios/...

.PHONY: test/android
test/android: ## Run Android tests only
	go test -v ./android/...

.PHONY: bench
bench: ## Run benchmark tests
	go test -bench=. -benchmem ./...

.PHONY: bench/vision
bench/vision: ## Run vision benchmarks
	go test -bench=. -benchmem ./internal/vision/

.PHONY: bench/locator
bench/locator: ## Run locator benchmarks
	go test -bench=. -benchmem ./internal/locator/

# =============================================================================
# Code Quality
# =============================================================================
.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

# =============================================================================
# Examples
# =============================================================================
.PHONY: example/basic
example/basic: ## Run basic example
	go run examples/basic_test.go

.PHONY: example/parallel
example/parallel: ## Run parallel example
	go run examples/parallel_test.go

# =============================================================================
# Documentation
# =============================================================================
.PHONY: doc
doc: ## Generate documentation
	@echo "Documentation available at: docs/plans/"

.PHONY: doc/server
doc/server: ## Start documentation server
	@echo "Open: docs/plans/2026-03-22-go-uop-design.md"

# =============================================================================
# CI/CD
# =============================================================================
.PHONY: ci
ci: tidy fmt vet lint test/cover ## Run full CI pipeline

.PHONY: release
release: ci build ## Build for release
	@echo "Build artifacts ready for release"

# =============================================================================
# Console (Web UI)
# =============================================================================
.PHONY: console/deps
console/deps: ## Install console frontend dependencies
	cd console && npm install

.PHONY: console/dev
console/dev: ## Run console frontend in dev mode (with proxy to :8080)
	cd console && npm run dev

.PHONY: console/build
console/build: ## Build console frontend
	cd console && npm run build

.PHONY: dev/console
dev/console: ## Run console backend in dev mode (serves on :8080, proxies to localhost:5173)
	go run ./cmd/console/main.go -dev -open

.PHONY: build/console
build/console: console/build ## Build console binary with embedded frontend
	go build -o bin/uop-console ./cmd/console/

.PHONY: run/console
run/console: ## Run console binary
	./bin/uop-console

# =============================================================================
# Clean
# =============================================================================
.PHONY: clean
clean: ## Clean build artifacts
	rm -rf coverage.out coverage.html
	rm -rf bin/uop-console
	rm -rf console/dist console/_out
	find . -name "*.test" -delete
	find . -name "*_mock.go" -delete
