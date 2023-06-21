.SILENT:

##
##> Dev

install: ## Install the project locally.
	brew upgrade
	brew install docker golang-migrate
	docker-compose pull
	docker-compose up -d --build
	make db-migration-run

start: ## Start the project.
	docker-compose start

stop: ## Stop the project.
	docker-compose stop

##
##> Database specific

db-migration-create: ## Create migrations files.
	migrate create -ext sql -dir db/migrations -seq ${MIGRATION_NAME}

db-migration-run: ## Run migrations.
	migrate -path db/migrations --database "postgresql://user:psw@localhost:5432/bank?sslmode=disable" up

.DEFAULT_GOAL := help
help:
	@grep -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'
.PHONY: help
