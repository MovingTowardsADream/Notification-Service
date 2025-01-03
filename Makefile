include .env

.PHONY: compose-up compose-down build migration-new-db docs_create

#APP_NAME := myapp
#BUILD_DIR := bin
#
#build:
#
#run: build
#
#clean:
#
#generate:
#
#docker-build:
#
#docker-run:



compose-up:
	docker-compose up --build -d postgres rabbitmq

compose-down:
	docker-compose down

build: compose-up
	docker-compose up --build

migration-new-db:
	go run ./cmd/migrator/main.go

gen-api:
	protoc -I ./api/proto \
      ./api/proto/notify/notify.proto \
      --go_out=./api/gen/go --go_opt=paths=source_relative \
      --go-grpc_out=./api/gen/go --go-grpc_opt=paths=source_relative

docs:
	protoc --doc_out=./docs --doc_opt=markdown,docs.md ./api/proto/notify/notify.proto

lint:
	golangci-lint run

test:
	go test ./... -v

all:

