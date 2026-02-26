-include .env
export

DB_DSN := host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable

.PHONY: sqlc migrate-up migrate-down run

sqlc:
	sqlc generate

migrate-up:
	goose -dir sql/migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir sql/migrations postgres "$(DB_DSN)" down

run:
	air
