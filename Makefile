ifeq ("$(wildcard .env)","")
# .env file does not exist
else
include .env
export
endif

PROJECT_NAME=auth_service

MIGRATE_CMD := $(shell which migrate)
MIGRATE_VERSION := v4.15.2

PROJECT_ROOT_PATH=$(CURDIR)
SERVICE_PATH_RELATIVE=cmd/auth/main.go
SERVICE_PATH=$(PROJECT_ROOT_PATH)/$(SERVICE_PATH_RELATIVE)
MIGRATION_PATH=$(PROJECT_ROOT_PATH)/migrations

SCRIPTS_PATH=$(PROJECT_ROOT_PATH)/scripts
CREATE_DEFAULT_ENV_SCRIPT=$(SCRIPTS_PATH)/create_default_env.sh
CREATE_LOCAL_ENV_SCRIPT=$(SCRIPTS_PATH)/create_local_env.sh
CREATE_TEST_ENV_SCRIPT=$(SCRIPTS_PATH)/create_test_env.sh

DEPLOYMENTS_PATH=$(PROJECT_ROOT_PATH)/deployments
DOCKER_COMPOSE_PATH=$(DEPLOYMENTS_PATH)/docker-compose.yaml
TEST_DOCKER_COMPOSE_PATH=$(DEPLOYMENTS_PATH)/docker-compose-test.yaml

DATABASE_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

MIGRATION_NAME := $(name)



.PHONY: env
env:
	$(CREATE_DEFAULT_ENV_SCRIPT)
	cp $(PROJECT_ROOT_PATH)/.env $(DEPLOYMENTS_PATH)/.env

.PHONY: local-env
local-env:
	$(CREATE_LOCAL_ENV_SCRIPT)
	cp $(PROJECT_ROOT_PATH)/.env $(DEPLOYMENTS_PATH)/.env

.PHONY: test-env
test-env:
	$(CREATE_TEST_ENV_SCRIPT)
	cp $(PROJECT_ROOT_PATH)/.env $(DEPLOYMENTS_PATH)/.env



.PHONY: docs
docs:
	swag init -g $(SERVICE_PATH_RELATIVE)

.PHONY: sqlc-gen
sqlc-gen:
	sqlc generate -f ./config/sqlc.yaml



.PHONY: build-service
build-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) build 

.PHONY: up-service
up-service: 
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) up -d

.PHONY: down-service
down-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) down

.PHONY: run-service
run-service:
	go run $(SERVICE_PATH)



.PHONY: build-test-service
build-test-service: 
	docker-compose -p $(PROJECT_NAME)_test -f $(TEST_DOCKER_COMPOSE_PATH) build

.PHONY: up-test-service
up-test-service: 
	docker-compose -p $(PROJECT_NAME)_test -f $(TEST_DOCKER_COMPOSE_PATH) up -d

.PHONY: down-test-service
down-test-service: 
	docker-compose -p $(PROJECT_NAME)_test -f $(TEST_DOCKER_COMPOSE_PATH) down

.PHONY: run-functional-tests
run-functional-tests:
	go test -tags=functional -count=1 ./tests/...



.PHONY: .check-migrate
.check-migrate:
ifeq ($(MIGRATE_CMD),)
	@echo "migrate not found, installing..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
endif

.PHONY: migration-up
migration-up: .check-migrate
	@echo "Running migrations up..."
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) up

.PHONY: migration-down
migration-down: .check-migrate
	@echo "Running migrations down..."
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) down

.PHONY: migration-create
migration-create: .check-migrate
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Migration name is required. Use: make migration-create name=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(MIGRATION_NAME)"
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(MIGRATION_NAME)



.PHONY: lines-count
lines-count:
	@echo 	Number of lines in GO files:
	@echo 	""[${shell find $(CURDIR) -name '*.go' -type f -print0 | xargs -0 cat | wc -l}]
