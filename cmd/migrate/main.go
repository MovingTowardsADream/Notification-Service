package main

import (
	"errors"
	"flag"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"Notification_Service/internal/infrastructure/config"
)

const (
	_defaultPathToMigrate = "./migrations"
)

func main() {
	var migratePath string

	flag.StringVar(&migratePath, "migrate-path", _defaultPathToMigrate, "path to migrations")
	flag.Parse()

	if migratePath == "" {
		panic("migrate-path is required")
	}

	cfg := config.MustLoad()

	if err := godotenv.Load(); err != nil {
		panic("failed reading env")
	}

	m, err := migrate.New(
		"file://"+migratePath,
		cfg.Storage.URL,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			// no migrations to apply
			return
		}
		panic(err)
	}
}
