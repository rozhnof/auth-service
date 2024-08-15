#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ENV_FILE_PATH=$SCRIPT_DIR"/../.env"

if [ ! -f $ENV_FILE_PATH ]; then
    POSTGRES_PASSWORD=$(openssl rand -base64 32 | tr -d '/+=')
    DATABASE_URL="postgres://user:$POSTGRES_PASSWORD@db:5432/mydb?sslmode=disable"
    SECRET_KEY=$(openssl rand -hex 32)

    echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD" > $ENV_FILE_PATH
    echo "DATABASE_URL=$DATABASE_URL" >> $ENV_FILE_PATH
    echo "SECRET_KEY=$SECRET_KEY" >> $ENV_FILE_PATH

    echo "secrets successfuly created"
fi
