# defining environment variables
PROJECT_PATH=$(shell dirname $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))
MIGRATIONS_DIR=db/migrations
name=default_migration_name

migrate-new: ## Create new migration
	dbmate --migrations-dir ${MIGRATIONS_DIR} new ${name}

migrate-create: ## Upgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} create ${name}

migrate-drop: ## Upgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} drop ${name}

migrate-up: ## Upgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} up

migrate-down: ## Downgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} down

cover:
	go test -cover ./...

mock-generate:
	go generate ./...
