# Auth Service

## Requirements

Before you begin, ensure you have the following installed:

- Docker
- Docker Compose
- Make 

## Installation and Running

### Step 1: Generate Secrets

```bash
make env
```

### Step 2: Build and Run with Docker Compose

```bash
docker compose up -d
```

## API Documentation

View the interactive API documentation on Swagger UI (the auth-service must be running):

[![Swagger UI](https://img.shields.io/badge/-API%20Documentation-blue)](https://localhost:8080/swagger/)