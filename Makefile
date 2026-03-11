-include .env
export

DB_DSN := host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable
GOOSE_MIGRATION_DIR := $(CURDIR)/sql/migrations

.PHONY: sqlc migrate-up migrate-down run stop

seed-admin:
	go run ./cmd/seed/main.go

sqlc:
	sqlc generate

migrate-up:
	goose up

migrate-down:
	goose down

## Run the application with Docker Compose and Air for live reloading
run:
	docker compose up -d
	air

## Stop the application and remove containers
stop:
	docker compose down
