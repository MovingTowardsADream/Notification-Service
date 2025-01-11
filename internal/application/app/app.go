package app

import (
	"context"

	"Notification_Service/internal/application/usecase"
	cltwilio "Notification_Service/internal/infrastructure/clients/twilio"
	"Notification_Service/internal/infrastructure/config"
	gwmessaging "Notification_Service/internal/infrastructure/gateway/messaging"
	"Notification_Service/internal/infrastructure/grpc"
	rmqclient "Notification_Service/internal/infrastructure/messaging/rabbitmq/client"
	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/infrastructure/smtp"
	amqprpc "Notification_Service/internal/infrastructure/workers/amqp_rpc"
	"Notification_Service/pkg/logger"
)

type App struct {
	Server          *grpc.Server
	MessagingServer *rmqserver.Server
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

	mesClient, err := rmqclient.New(
		cfg.Messaging.URL,
		cfg.Messaging.Server.RPCExchange,
		cfg.Messaging.Client.RPCExchange,
		cfg.Messaging.Topics,
		rmqclient.ConnAttempts(cfg.Messaging.Client.Attempts),
		rmqclient.ConnWaitTime(cfg.Messaging.Client.WaitTime),
		rmqclient.Timeout(cfg.Messaging.Client.Timeout),
	)

	if err != nil {
		panic("messaging client connection error" + err.Error())
	}

	usersRepo := postgres.NewUsersRepo(storage)

	gateway := gwmessaging.NewNotifyGateway(mesClient)

	notifySender := usecase.NewNotifySender(l, usersRepo, gateway)
	editInfo := usecase.NewEditInfo(l, usersRepo)

	gRPCServer := grpc.New(
		l, notifySender, editInfo, grpc.Port(cfg.GRPC.Port),
	)

	phoneSenderClient := cltwilio.NewClient(cfg.PhoneSender.AccountSID, cfg.PhoneSender.AuthToken, cfg.PhoneSender.MessagingServiceSID)

	senderPhone := cltwilio.NewWorkerPhone(phoneSenderClient)

	smtpClient := smtp.New(
		smtp.Params{
			Domain:   cfg.Domain,
			Username: cfg.UserName,
			Password: cfg.Password,
		},
		smtp.Port(cfg.SMTP.Port),
	)

	mailWorker := smtp.NewWorkerMail(smtpClient)

	rmqRouter := amqprpc.NewRouter(mailWorker, senderPhone)

	mesServer, err := rmqserver.New(
		cfg.Messaging.URL,
		cfg.Messaging.Server.RPCExchange,
		cfg.Messaging.Topics,
		rmqRouter,
		l,
		rmqserver.GoroutinesCount(cfg.Messaging.Server.GoroutinesCount),
		rmqserver.ConnAttempts(cfg.Messaging.Server.Attempts),
		rmqserver.ConnWaitTime(cfg.Messaging.Server.WaitTime),
		rmqserver.Timeout(cfg.Messaging.Server.Timeout),
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
