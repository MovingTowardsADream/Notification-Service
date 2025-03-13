include .env

NAMESPACE = default

.PHONY: compose-up
compose-up:
	docker compose up --build -d postgres rabbitmq loki promtail grafana

.PHONY: build
build: compose-up
	docker compose up --build -d

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: migrate
migrate:
	go run ./cmd/migrate/main.go

.PHONY: gen-api
gen-api:
	protoc -I ./api/proto \
		-I ./api/proto/validate \
		-I ./api/proto \
        ./api/proto/notify/notify.proto \
        --go_out=./api/gen/go --go_opt=paths=source_relative \
        --go-grpc_out=./api/gen/go --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=./api/gen/go --grpc-gateway_opt=paths=source_relative \
        --validate_out="lang=go:./api/gen/go"

.PHONY: docs
docs:
	protoc -I ./api/proto \
           -I ./api/proto/validate \
           -I ./api/proto/google \
           ./api/proto/notify/notify.proto \
           --openapiv2_out=./docs/openapi --openapiv2_opt=logtostderr=true \
           --doc_out=./docs/protocol --doc_opt=html,documentation.html

test:
	go test ./... -v

lint:
	golangci-lint run

