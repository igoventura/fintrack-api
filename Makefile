# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: migrate rollback tidy test help

help:
	@echo "Available commands:"
	@echo "  make migrate   - Run database migrations using tern"
	@echo "  make rollback  - Rollback the last migration using tern"
	@echo "  make tidy      - Run go mod tidy"
	@echo "  make test      - Run all tests"

migrate:
	@echo "Running migrations..."
	@tern migrate -m migrations/

rollback:
	@echo "Rolling back last migration..."
	@tern rollback -m migrations/

tidy:
	@go mod tidy

test:
	@echo "Running tests..."
	@go test ./...
