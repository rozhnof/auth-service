FROM golang:1.22-alpine

WORKDIR /app

RUN apk add --no-cache bash git g++ make openssl
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

EXPOSE 8080