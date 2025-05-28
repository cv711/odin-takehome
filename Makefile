# Default target
.DEFAULT_GOAL := help

# Help documentation
help:
	@echo "Available targets:"
	@echo "server      - Run Go server"
	@echo "db-generate - Generate SQL code"
	@echo "db-migrate  - Run database migrations up"
	@echo "db-status   - Check migration status"

# Server target
server:
	go run ./server

# Database migration targets
.PHONY: _goose db-generate db-migrate db-status
define _goose
	GOOSE_MIGRATION_DIR=server/db/migrations goose postgres \
	"user=odin dbname=odinexercise sslmode=disable host=127.0.0.1 port=5432 password=exercise" $(1)
endef

db-generate:
	sqlc generate

db-migrate:
	@$(call _goose,up)

db-status:
	@$(call _goose,status)