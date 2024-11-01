ifeq ("$(wildcard .env)","")
# .env file does not exist
else
include .env
export
endif


PROJECT_NAME=auth_service

MIGRATE_CMD := $(shell which migrate)
MIGRATE_VERSION := v4.15.2

TEST_CMD := $(shell which gotestsum)

PROJECT_ROOT_PATH=$(CURDIR)
SERVICE_PATH_RELATIVE=cmd/http/main.go
SERVICE_PATH=$(PROJECT_ROOT_PATH)/$(SERVICE_PATH_RELATIVE)
MIGRATION_PATH=$(PROJECT_ROOT_PATH)/migrations

SCRIPTS_PATH=$(PROJECT_ROOT_PATH)/scripts
CREATE_DEFAULT_ENV_SCRIPT=$(SCRIPTS_PATH)/create_default_env.sh
CREATE_LOCAL_ENV_SCRIPT=$(SCRIPTS_PATH)/create_local_env.sh
EXPORT_TEST_ENV_SCRIPT=$(SCRIPTS_PATH)/export_test_env.sh

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

.PHONY: docs
docs:
	swag init -g $(SERVICE_PATH_RELATIVE)

build-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) build 

up-service: 
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) up -d

down-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) down

run-service:
	go run $(SERVICE_PATH)



run-integration-tests:
	. $(EXPORT_TEST_ENV_SCRIPT) && docker-compose -p $(PROJECT_NAME)_test -f $(TEST_DOCKER_COMPOSE_PATH) up -d --build postgres migrate
	. $(EXPORT_TEST_ENV_SCRIPT) && go test -tags=integration -count=1 ./tests/...
	. $(EXPORT_TEST_ENV_SCRIPT) && docker-compose -p $(PROJECT_NAME)_test -f $(TEST_DOCKER_COMPOSE_PATH) down -v

run-unit-tests:
	go test -count=1 ./internal/...



.check-migrate:
ifeq ($(MIGRATE_CMD),)
	@echo "migrate not found, installing..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
endif

migration-up: .check-migrate
	@echo "Running migrations up..."
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) up

migration-down: .check-migrate
	@echo "Running migrations down..."
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) down

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
