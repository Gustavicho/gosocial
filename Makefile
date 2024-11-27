include .env
MIGRATIONS_DIR = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database=$(DB_DSN) up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database=$(DB_DSN) down $(filter-out $@,$(MAKECMDGOALS))