#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_ROOT_PATH=$SCRIPT_DIR"/.."
ENV_FILE_PATH=$PROJECT_ROOT_PATH"/.env"

CONFIG_PATH=config/config.yaml
POSTGRES_ADDRESS=postgres
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=auth_db
POSTGRES_SSLMODE=disable

REDIS_ADDRESS=redis
REDIS_PORT=6379 
REDIS_USER=user 
REDIS_PASSWORD=password 
REDIS_DB=0

SECRET_KEY=$2a$10$Ck6JA2dOt/DdCS28eGuah.PKp/.1BI4sebcFgPKNJEXdVqAQ4S6KC
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

echo "CONFIG_PATH=$CONFIG_PATH" > $ENV_FILE_PATH
echo "POSTGRES_ADDRESS=$POSTGRES_ADDRESS" >> $ENV_FILE_PATH
echo "POSTGRES_PORT=$POSTGRES_PORT" >> $ENV_FILE_PATH
echo "POSTGRES_USER=$POSTGRES_USER" >> $ENV_FILE_PATH
echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD" >> $ENV_FILE_PATH
echo "POSTGRES_DB=$POSTGRES_DB" >> $ENV_FILE_PATH
echo "POSTGRES_SSLMODE=$POSTGRES_SSLMODE" >> $ENV_FILE_PATH

echo "REDIS_ADDRESS=$REDIS_ADDRESS" >>  $ENV_FILE_PATH
echo "REDIS_PORT=$REDIS_PORT" >>  $ENV_FILE_PATH
echo "REDIS_USER=$REDIS_USER" >>  $ENV_FILE_PATH
echo "REDIS_PASSWORD=$REDIS_PASSWORD" >>  $ENV_FILE_PATH
echo "REDIS_DB=$REDIS_DB" >>  $ENV_FILE_PATH

echo "SECRET_KEY=$SECRET_KEY" >> $ENV_FILE_PATH
echo "GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID" >> $ENV_FILE_PATH
echo "GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET" >> $ENV_FILE_PATH