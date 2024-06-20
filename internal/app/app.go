package app

import (
	"Notification_Service/configs"
	gateway_rabbitmq "Notification_Service/internal/notify/gateway/rabbitmq"
	grpcserver2 "Notification_Service/internal/notify/grpc/grpcserver"
	"Notification_Service/internal/notify/repository/postgresdb"
	"Notification_Service/internal/notify/usecase/edit_info"
	"Notification_Service/internal/notify/usecase/send_notify"
	"Notification_Service/pkg/postgres"
	rmq_client "Notification_Service/pkg/rabbitmq/client"
	rmq_server "Notification_Service/pkg/rabbitmq/server"
	"log/slog"
)

type App struct {
	Server    *grpcserver2.Server
	RMQServer *rmq_server.Server
	DB        *postgres.Postgres
}

func New(l *slog.Logger, cfg *configs.Config) *App {

	// Connect postgres db
	pg, err := postgres.NewPostgresDB(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		panic("app - New - postgres.NewPostgresDB: " + err.Error())
	}

	rmqClient, err := rmq_client.NewRabbitMQClient(cfg.RMQ.URL, cfg.RMQ.ServerExchange, cfg.RMQ.ClientExchange)
	if err != nil {
		panic("app - Run - rmqServer - server.New" + err.Error())
	}

	notifyRepo := postgresdb.NewNotifyRepo(pg)

	gateway := gateway_rabbitmq.New(rmqClient)

	notifySend := send_notify.New(l, notifyRepo, gateway)
	editInfo := edit_info.New(l, notifyRepo)

	gRPCServer := grpcserver2.New(l, notifySend, editInfo, grpcserver2.Port(cfg.Port))

	return &App{
		Server: gRPCServer,
		DB:     pg,
	}
}
