package gotests

import (
	"context"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"

	"Notification_Service/internal/infrastructure/config"
)

const (
	_defaultStorageImage    = "postgres:15"
	_defaultMessagingImaged = "rabbitmq:3-management"
)

func runDatabaseIntegration(ctx context.Context, cfg *config.Config) func() {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(_defaultStorageImage),
		postgres.WithDatabase(cfg.Storage.DBName),
		postgres.WithUsername(cfg.Storage.Username),
		postgres.WithPassword(cfg.Storage.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(1),
		),
	)
	if err != nil {
		panic(err)
	}

	cancel := func() {
		_ = pgContainer.Terminate(ctx)
	}

	return cancel
}

func runMessagingIntegration(ctx context.Context, cfg *config.Config) func() {
	rabbitContainer, err := rabbitmq.RunContainer(ctx,
		testcontainers.WithImage(_defaultMessagingImaged),
		rabbitmq.WithAdminUsername(cfg.Messaging.Username),
		rabbitmq.WithAdminPassword(cfg.Messaging.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("server startup complete").WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		panic(err)
	}

	cancel := func() {
		_ = rabbitContainer.Terminate(ctx)
	}

	return cancel
}
