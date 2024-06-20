include .env

.PHONY:

compose-up:
	docker-compose up --build -d postgres rabbitmq

compose-down:
	docker-compose down

build: compose-up
	docker-compose up --build

migration-new-db:
	go run ./cmd/migrator/main.go

docs_create:
	protoc --doc_out=./docs --doc_opt=markdown,docs.md ./protos/proto/notify/notify.proto
