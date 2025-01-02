package app

import (
	"context"

	"Notification_Service/internal/application/usecase"
	"Notification_Service/internal/infrastructure/config"
	gateway_messaging "Notification_Service/internal/infrastructure/gateway/messaging"
	"Notification_Service/internal/infrastructure/grpc"
	rmq_client "Notification_Service/internal/infrastructure/messaging/rabbitmq/client"
	rmq_server "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/infrastructure/smtp"
	"Notification_Service/internal/infrastructure/workers/amqp_rpc"
	"Notification_Service/pkg/logger"
)

type App struct {
	Server          *grpc.Server
	MessagingServer *rmq_server.Server
	Storage         *postgres.Postgres
}

func New(ctx context.Context, l *logger.Logger, cfg *config.Config) *App {
	storage, err := postgres.New(
		ctx,
		cfg.Storage.URL,
		postgres.MaxPoolSize(cfg.Storage.PoolMax),
		postgres.ConnAttempts(cfg.Storage.ConnAttempts),
		postgres.ConnTimeout(cfg.Storage.ConnTimeout),
	)

	if err != nil {
		panic("storage connection error: " + err.Error())
	}

	err = storage.Ping(ctx)

	if err != nil {
		panic("storage ping error: " + err.Error())
	}

	mesClient, err := rmq_client.New(
		cfg.Messaging.URL,
		cfg.Messaging.Server.RPCExchange,
		cfg.Messaging.Client.RPCExchange,
		rmq_client.ConnAttempts(cfg.Messaging.Client.Attempts),
		rmq_client.ConnWaitTime(cfg.Messaging.Client.WaitTime),
		rmq_client.Timeout(cfg.Messaging.Client.Timeout),
	)

	if err != nil {
		panic("messaging client connection error" + err.Error())
	}

	usersRepo := postgres.NewUsersRepo(storage)

	gateway := gateway_messaging.NewNotifyGateway(mesClient)

	notifySend := usecase.NewNotifySend(l, usersRepo, gateway)
	editInfo := usecase.NewEditInfo(l, usersRepo)

	gRPCServer := grpc.New(
		l, notifySend, editInfo, grpc.Port(cfg.GRPC.Port),
	)

	SMTPClient := smtp.New(
		&smtp.Params{
			Domain:   cfg.Domain,
			Username: cfg.UserName,
			Password: cfg.Password,
			Mail:     cfg.Notify.Mail,
		},
		smtp.Port(cfg.SMTP.Port),
	)

	workerUseCase := smtp.NewNotifyWorker(SMTPClient)

	rmqRouter := amqp_rpc.NewRouter(workerUseCase)

	rmqServer, err := rmq_server.New(
		cfg.Messaging.URL,
		cfg.Messaging.Server.RPCExchange,
		rmqRouter,
		l,
		rmq_server.GoroutinesCount(cfg.Messaging.Server.GoroutinesCount),
		rmq_server.ConnAttempts(cfg.Messaging.Server.Attempts),
		rmq_server.ConnWaitTime(cfg.Messaging.Server.WaitTime),
		rmq_server.Timeout(cfg.Messaging.Server.Timeout),
	)

	if err != nil {
		panic("messaging server connection error" + err.Error())
	}

	return &App{
		Server:          gRPCServer,
		MessagingServer: rmqServer,
		Storage:         storage,
	}
}
