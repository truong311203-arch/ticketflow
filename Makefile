.PHONY: all build run test test-integration test-coverage clean tidy help

APP_NAME := ticketflow-api
CMD_DIR := ./cmd/api
BUILD_DIR := ./bin

ifeq ($(OS),Windows_NT)
	BINARY := $(BUILD_DIR)/$(APP_NAME).exe
else
	BINARY := $(BUILD_DIR)/$(APP_NAME)
endif

all: test build

build: ## Build the API binary
	@echo Building $(BINARY)...
	go build -v -o $(BINARY) $(CMD_DIR)

run: ## Run the API server locally
	go run $(CMD_DIR)

test: ## Run unit tests
	go test -v -race ./internal/...

test-integration: ## Run integration tests
	go test -v -race ./tests/integration/...

test-coverage: ## Run all tests with coverage report
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

tidy: ## Tidy and verify Go modules
	go mod tidy
	go mod verify

clean: ## Clean build artifacts and coverage files
	@if exist bin rmdir /s /q bin 2>nul || rm -rf bin 2>/dev/null || true
	@del /f /q coverage.out 2>nul || rm -f coverage.out 2>/dev/null || true

help: ## Show this help message
	@echo Usage: make [target]
	@echo Targets:
	@echo   build            Build the API binary
	@echo   run              Run the API server locally
	@echo   test             Run internal unit tests
	@echo   test-integration Run integration tests
	@echo   test-coverage    Run tests with code coverage summary
	@echo   tidy             Tidy and verify dependencies
	@echo   clean            Remove build and coverage artifacts
