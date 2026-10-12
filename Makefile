# Variables
BINARY_NAME=forum
GOFLAGS=-v
GO=go

# DB Container Variables
DB_CONTAINER_NAME=forumdb
DB_USER=godbless
DB_NAME=forumdb

# Build variables
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

.PHONY: all build clean test lint fmt run dev db-up db-down db-terminal db-logs install help

all: fmt test build

## build: Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/forum-app

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

## test: Run all tests
test:
	@echo "Running tests..."
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

## lint: Run linter
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	$(GO) fmt ./...
	$(GO) vet ./...

## run: Build and run the application
run: build
	./bin/$(BINARY_NAME)

## dev: Build and run reading .env from disk (if present)
dev: build
	@test -f .env || cp .env.example .env; \
	./bin/$(BINARY_NAME)

## Migrate DB Schema to Database
db-migrate-up:
db-migrate-down:

## db-terminal: Bring up psql interactive terminal
db-terminal:
	@echo "Waiting for Postgres to accept connections..."; \
	until docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) >/dev/null 2>&1; do sleep 1; done; \
	echo "Postgres ready"


## db-logs: Tail Postgres logs
db-logs:
	docker logs $(DB_CONTAINER_NAME)

## install: Download dependencies
install:
	$(GO) mod download
	$(GO) mod tidy

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'