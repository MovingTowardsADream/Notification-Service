package app

import (
	"Notification_Service/configs"
	grpcserver2 "Notification_Service/internal/notify/grpc/grpcserver"
	"Notification_Service/internal/notify/repository/postgresdb"
	"Notification_Service/internal/notify/usecase/edit_info"
	"Notification_Service/internal/notify/usecase/send_notify"
	"Notification_Service/pkg/postgres"
	"log/slog"
)

type App struct {
	Server *grpcserver2.Server
	DB     *postgres.Postgres
}

func New(l *slog.Logger, cfg *configs.Config) *App {

	// Connect postgres db
	pg, err := postgres.NewPostgresDB(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		panic("app - New - postgres.NewPostgresDB: " + err.Error())
	}

	notifyRepo := postgresdb.NewNotifyRepo(pg)

	notifySend := send_notify.New(l, notifyRepo)
	editInfo := edit_info.New(l, notifyRepo)

	gRPCServer := grpcserver2.New(l, notifySend, editInfo, grpcserver2.Port(cfg.Port))

	return &App{
		Server: gRPCServer,
		DB:     pg,
	}
}
