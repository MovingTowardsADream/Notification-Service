package app

import (
	"Notification_Service/configs"
	gateway_rabbitmq "Notification_Service/internal/notify/gateway/rabbitmq"
	grpcserver2 "Notification_Service/internal/notify/grpc/grpcserver"
	"Notification_Service/internal/notify/repository/postgresdb"
	"Notification_Service/internal/notify/usecase/edit_info"
	"Notification_Service/internal/notify/usecase/send_notify"
	"Notification_Service/internal/notifyWorkers/controller/amqp_rpc"
	"Notification_Service/internal/notifyWorkers/usecase"
	"Notification_Service/pkg/postgres"
	rmq_client "Notification_Service/pkg/rabbitmq/client"
	rmq_server "Notification_Service/pkg/rabbitmq/server"
	"log/slog"
)

const _defaultPathToMigrate = "./schema/0000001_init.up.sql"

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

	// Migrate database schema
	err = pg.Migrate(_defaultPathToMigrate)
	if err != nil {
		panic("app - Run - pg.Migrate: " + err.Error())
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

	workerUseCase := usecase.NewNotifyWorker()

	// Init rabbitMQ RPC Server
	rmqRouter := amqp_rpc.NewRouter(workerUseCase)

	rmqServer, err := rmq_server.New(
		cfg.RMQ.URL,
		cfg.RMQ.ServerExchange,
		rmqRouter,
		l,
		rmq_server.DefaultGoroutinesCount(cfg.App.CountWorkers),
	)
	if err != nil {
		panic("app - Run - rmqServer - server.New" + err.Error())
	}

	return &App{
		Server:    gRPCServer,
		RMQServer: rmqServer,
		DB:        pg,
	}
}
