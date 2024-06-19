package app

import (
	"Notification_Service/configs"
	grpcserver2 "Notification_Service/internal/notify/grpc/grpcserver"
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

	gRPCServer := grpcserver2.New(l, grpcserver2.Port(cfg.Port))

	return &App{
		Server: gRPCServer,
		DB:     pg,
	}
}
