# Auth Service

## Requirements

- Docker
- Docker Compose
- Make 

## Service

### Build and run service

```bash
make env
make build-service
make up-service
```

### Stop service

```bash
make down-service
```

## Tests

### Functional tests

#### Build and run service with test environment

```bash
make test-env
make build-test-service
make up-test-service
```

#### Run tests

```bash
make run-functional-tests
```

#### Stop service

```bash
make down-test-service
```

## Database

### Up migrations

```bash
make migration-up
```

### Down migrations

```bash
make migration-down
```

### Create new migration

```bash
make migration-create name=<migration_name>
```

## API

[![Swagger UI](https://img.shields.io/badge/-API%20Documentation-blue)](http://localhost:8080/swagger/index.html)