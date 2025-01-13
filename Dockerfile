FROM alpine:3.15 as root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

FROM golang:1.22 as builder
WORKDIR /notification-service
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o notify ./cmd/notify/main.go

FROM scratch as final

MAINTAINER Matvet Smirnov "MovingTowardsADream"

LABEL maintainer="MovingTowardsADream"
LABEL version="v1.0.0"
LABEL description="This is notification service"

ARG defaultPort=8080

COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --chown=1001:1001 --from=builder /notification-service /notification-service
COPY --chown=1001:1001 ./configs ./configs
COPY --chown=1001:1001 ./.env .

USER app
EXPOSE $defaultPort

ENTRYPOINT ["./notification-service/notify"]
