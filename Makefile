.SILENT:

##
##> Dev

install: ## Install the project locally.
	brew upgrade
	brew install docker golang-migrate sqlc go
	docker-compose pull
	docker-compose up -d --build
	make db-migration-run
	go mod download

start: ## Start the project.
	docker-compose start

stop: ## Stop the project.
	docker-compose stop

##
##> Database specific

db-migration-create: ## Create DB migrations files.
	migrate create -ext sql -dir db/migrations -seq ${MIGRATION_NAME}

db-migration-run: ## Run DB migrations.
	migrate -path db/migrations --database "postgresql://user:psw@localhost:5432/bank?sslmode=disable" up

db-orm-generate: ## Generate CRUD ORM according to DB.
	sqlc generate

##
##> Test specific
test-run:
	go test -cover ./...

.DEFAULT_GOAL := help
help:
	@grep -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'
.PHONY: help
