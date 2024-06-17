package app

import (
	"Notification_Service/configs"
	"Notification_Service/pkg/postgres"
	"log/slog"
)

type App struct {
	log *slog.Logger

	DB *postgres.Postgres
}

func New(l *slog.Logger, cfg *configs.Config) *App {

	// Connect postgres db
	pg, err := postgres.NewPostgresDB(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		panic("app - New - postgres.NewPostgresDB: " + err.Error())
	}

	return &App{
		log: l,
		DB:  pg,
	}
}
