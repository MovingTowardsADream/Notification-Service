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
#
# dounload validate
# curl -O https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/main/validate/validate.proto

# и ген и валид и дока
#protoc -I ./api/proto \
#       -I ./api/proto/validate \
#       -I ./api/proto/google \
#       ./api/proto/notify/notify.proto \
#       --go_out=./api/gen/go --go_opt=paths=source_relative \
#       --go-grpc_out=./api/gen/go --go-grpc_opt=paths=source_relative \
#       --validate_out="lang=go:./api/gen/go" \
#       --openapiv2_out=./docs/openapi --openapiv2_opt=logtostderr=true


compose-up:
	docker-compose up --build -d postgres rabbitmq

compose-down:
	docker-compose down

build: compose-up
	docker-compose up --build

migrate:
	go run ./cmd/migrate/main.go

gen-api:
	protoc -I ./api/proto -I ./api/proto/validate \
    	./api/proto/notify/notify.proto \
    	--go_out=./api/gen/go --go_opt=paths=source_relative \
    	--go-grpc_out=./api/gen/go --go-grpc_opt=paths=source_relative \
    	--validate_out="lang=go:./api/gen/go"

docs:
	protoc -I ./api/proto \
           -I ./api/proto/validate \
           -I ./api/proto/google \
           ./api/proto/notify/notify.proto \
           --openapiv2_out=./docs/openapi --openapiv2_opt=logtostderr=true \
           --doc_out=./docs/protocol --doc_opt=html,documentation.html

lint:
	golangci-lint run

test:
	go test ./... -v

all:

