package app

import (
	"context"

	"Notification_Service/internal/application/usecase"
	clTwilio "Notification_Service/internal/infrastructure/clients/twilio"
	"Notification_Service/internal/infrastructure/config"
	gatewayMessaging "Notification_Service/internal/infrastructure/gateway/messaging"
	"Notification_Service/internal/infrastructure/grpc"
	rmqClient "Notification_Service/internal/infrastructure/messaging/rabbitmq/client"
	rmqServer "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/infrastructure/smtp"
	"Notification_Service/internal/infrastructure/workers/amqp_rpc"
	"Notification_Service/pkg/logger"
)

type App struct {
	Server          *grpc.Server
	MessagingServer *rmqServer.Server
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

	mesClient, err := rmqClient.New(
		cfg.Messaging.URL,
		cfg.Messaging.Server.RPCExchange,
		cfg.Messaging.Client.RPCExchange,
		cfg.Messaging.Topics,
		rmqClient.ConnAttempts(cfg.Messaging.Client.Attempts),
		rmqClient.ConnWaitTime(cfg.Messaging.Client.WaitTime),
		rmqClient.Timeout(cfg.Messaging.Client.Timeout),
	)

	if err != nil {
		panic("messaging client connection error" + err.Error())
	}

	usersRepo := postgres.NewUsersRepo(storage)

	gateway := gatewayMessaging.NewNotifyGateway(mesClient)

	notifySend := usecase.NewNotifySend(l, usersRepo, gateway)
	editInfo := usecase.NewEditInfo(l, usersRepo)

	gRPCServer := grpc.New(
		l, notifySend, editInfo, grpc.Port(cfg.GRPC.Port),
	)

	phoneSenderClient := clTwilio.NewClient(cfg.PhoneSender.AccountSID, cfg.PhoneSender.AuthToken, cfg.PhoneSender.MessagingServiceSID)

	_ = phoneSenderClient

	smtpClient := smtp.New(
		smtp.Params{
			Domain:   cfg.Domain,
			Username: cfg.UserName,
			Password: cfg.Password,
		},
		smtp.Port(cfg.SMTP.Port),
	)

	workerUseCase := smtp.NewNotifyWorker(smtpClient)

	rmqRouter := amqp_rpc.NewRouter(workerUseCase)

	mesServer, err := rmqServer.New(
		cfg.Messaging.URL,
		cfg.Messaging.Server.RPCExchange,
		cfg.Messaging.Topics,
		rmqRouter,
		l,
		rmqServer.GoroutinesCount(cfg.Messaging.Server.GoroutinesCount),
		rmqServer.ConnAttempts(cfg.Messaging.Server.Attempts),
		rmqServer.ConnWaitTime(cfg.Messaging.Server.WaitTime),
		rmqServer.Timeout(cfg.Messaging.Server.Timeout),
	)

	if err != nil {
		panic("messaging server connection error" + err.Error())
	}

	return &App{
		Server:          gRPCServer,
		MessagingServer: mesServer,
		Storage:         storage,
	}
}
