MIGRATION_FOLDER=$(CURDIR)/migrations
MAIN_PATH=$(CURDIR)/cmd/main.go
DB_DRIVER=postgres
SCRIPTS_FOLDER=$(CURDIR)/scripts

migration-create:
	@if [ -z "$(name)" ]; then \
	  echo "Error: migration name not provided"; \
	  echo "Usage: make create-migration name='migration-name'"; \
	  exit 1; \
	fi
	goose -dir $(MIGRATION_FOLDER) create $(name) sql

.PHONY: env
env:
	$(SCRIPTS_FOLDER)/generate_env.sh

.PHONY: generate-certificate
generate-certificate:
	$(SCRIPTS_FOLDER)/generate_certificate.sh

.PHONY: run-service
run-service: 
	go run $(MAIN_PATH)

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" $(DB_DRIVER) "$(DATABASE_URL)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" $(DB_DRIVER) "$(DATABASE_URL)" down


