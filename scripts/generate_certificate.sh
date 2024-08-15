#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
CERT_FILE=$SCRIPT_DIR/../cert.pem
KEY_FILE=$SCRIPT_DIR/../key.pem

if [ ! -f $CERT_FILE ]; then
    echo "Certificate does not exist. Generating..."
    echo | openssl genrsa -out "$KEY_FILE" 4096
    echo | openssl req -new -key "$KEY_FILE" -out "$CERT_FILE".csr -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com"
    openssl x509 -req -days 365 -in "$CERT_FILE".csr -signkey "$KEY_FILE" -out "$CERT_FILE"
    rm "$CERT_FILE".csr
    echo "Certificate has been successfully generated."
fi

if [ ! -f $KEY_FILE ]; then
    echo "Key does not exist. Generating..."
    echo | openssl genrsa -out "$KEY_FILE" 4096
    echo "Key has been successfully generated."
fi
