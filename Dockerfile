
FROM golang:1.22-alpine AS builder
WORKDIR /notification-service
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/app/main.go

FROM alpine:latest
WORKDIR /notification-service
COPY --from=builder /notification-service ./

EXPOSE 8080

CMD ["./notification-service"]