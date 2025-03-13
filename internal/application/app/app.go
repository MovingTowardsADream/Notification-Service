package app

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"Notification_Service/internal/application/usecase"
	"Notification_Service/internal/infrastructure/clients/smtp"
	cltwilio "Notification_Service/internal/infrastructure/clients/twilio"
	"Notification_Service/internal/infrastructure/config"
	gwmessaging "Notification_Service/internal/infrastructure/gateway/messaging"
	"Notification_Service/internal/infrastructure/grpc"
	rmqclient "Notification_Service/internal/infrastructure/messaging/rabbitmq/client"
	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/infrastructure/observ/metrics"
	"Notification_Service/internal/infrastructure/observ/trace"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/infrastructure/repository/postgres/users"
	amqprpc "Notification_Service/internal/infrastructure/workers/amqp_rpc"
	"Notification_Service/pkg/hasher"
	"Notification_Service/pkg/logger"
	"Notification_Service/pkg/utils"
)

type App struct {
	Server          *grpc.Server
	MessagingServer *rmqserver.Server
	Storage         *postgres.Postgres
	Tracer          *trace.TracesProvider
	MetricsServer   *metrics.Server
}

func New(ctx context.Context, l logger.Logger, cfg *config.Config) *App {
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

	tracer, err := trace.New(
		ctx,
		utils.FormatAddress(cfg.Observability.Trace.Host, cfg.Observability.Trace.Port),
		cfg.App.Name,
		trace.Enabled(cfg.Observability.Trace.Enabled),
		trace.InitialInterval(cfg.Observability.Trace.InitialInterval),
		trace.MaxInterval(cfg.Observability.Trace.MaxInterval),
		trace.MaxElapsedTime(cfg.Observability.Trace.MaxElapsedTime),
	)

	if err != nil {
		panic("trace provider connection error" + err.Error())
	}

	usersRepo := users.NewUsersRepo(storage)

	gateway := gwmessaging.NewNotifyGateway(mesClient)

	hash := hasher.NewSHA1Hasher(cfg.Security.PasswordSalt)

	notifySender := usecase.NewNotifySender(l, usersRepo, gateway)
	editInfo := usecase.NewUserInfo(l, usersRepo, hash)

	reg := prometheus.NewRegistry()
	m := metrics.New(reg, cfg.App.Name)

	gRPCServer := grpc.New(
		l, m, notifySender, editInfo, grpc.Port(cfg.GRPC.Port),
	)

	metricsServer := metrics.NewMetricsServer(
		utils.FormatAddress("", cfg.Observability.Metrics.Port),
		reg,
		l,
	)

	phoneSenderClient := cltwilio.NewClient(cfg.PhoneSender.AccountSID, cfg.PhoneSender.AuthToken, cfg.PhoneSender.MessagingServiceSID)

	senderPhone := cltwilio.NewWorkerPhone(phoneSenderClient)

	smtpClient := smtp.New(
		smtp.Params{
			Domain:   cfg.SMTP.Domain,
			Username: cfg.SMTP.UserName,
			Password: cfg.SMTP.Password,
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
		Tracer:          tracer,
		MetricsServer:   metricsServer,
	}
}
