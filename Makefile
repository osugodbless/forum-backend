# Load .env file automatically if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Path to the database migration script
MIGRATE_SCRIPT := ./scripts/migrate.sh

# Variables
BINARY_NAME=forum
GOFLAGS=-v
GO=go

# DB Container Variables
DB_CONTAINER_NAME=forumdb
DB_USER=forum
DB_NAME=forumdb

# Build variables
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

.PHONY: all build clean test lint fmt run dev check-env migrate-up migrate-down migrate-status db-terminal db-logs install help

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

## check-env: Verify that DATABASE_URL is available
check-env:
	@if [ -z "$$DB_URL" ]; then \
		echo "Error: DATABASE_URL is not set. Define it in .env or pass it inline (e.g. make migrate-up DATABASE_URL=...)"; \
		exit 1; \
	fi

## migrate-up: Run all migrations up (auth -> profile -> content)
migrate-up: check-env
	@DATABASE_URL="$$DB_URL" $(MIGRATE_SCRIPT) up

## migrate-down: Roll back the last migration step in reverse order (content -> profile -> auth)
migrate-down: check-env
	@DATABASE_URL="$$DB_URL" $(MIGRATE_SCRIPT) down

## migrate-status: Check the current migration version for each module
migrate-status: check-env
	@DATABASE_URL="$$DB_URL" $(MIGRATE_SCRIPT) status

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