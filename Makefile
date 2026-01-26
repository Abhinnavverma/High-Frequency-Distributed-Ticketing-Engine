# Makefile

# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Construct the DB URL using the variables from .env
# This keeps your secrets in one place (.env) and your logic in another (Makefile)
DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: docker-up docker-down migrate-up migrate-down migrate-create

# Start the containers in the background
docker-up:
	docker-compose up -d

# Stop the containers
docker-down:
	docker-compose down

# Apply all "up" migrations
migrate-up:
	migrate -path internals/db/migrations -database "$(DB_URL)" up

# Revert the last migration (Safety switch)
migrate-down:
	migrate -path internals/db/migrations -database "$(DB_URL)" down 1

# Force a specific version (useful if migration gets stuck)
migrate-force:
	migrate -path internals/db/migrations -database "$(DB_URL)" force $(version)

# Create a new migration file pair
# Usage: make migrate-create name=add_users_table
migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)
run-dev:
	go run cmd/Server/main.go