MIGRATIONS_DIR := ./internal/infrastructure/persistence/migrations
POSTGRES_DSN   ?= postgresql://postgres:root@localhost:5432/users?sslmode=disable

-include .env
export

.PHONY: setup swagger migrate-up migrate-down run

setup:
	@cp -n .env.example .env 2>/dev/null && echo ".env created from .env.example" || echo ".env already exists"
	@go mod tidy
	@go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	@swag init -g cmd/api/main.go --output docs

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" up

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" down

run:
	@go run ./cmd/api
