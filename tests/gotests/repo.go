package gotests

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"

	// package for initializing migration on PostgreSQL.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"Notification_Service/internal/infrastructure/config"
	rmqclient "Notification_Service/internal/infrastructure/messaging/rabbitmq/client"
	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/pkg/logger"
)

const pathToMigrate = "../../../migrations"

var ErrRepositoryImpl = errors.New("errors repository implementation")

type Repository interface {
	UnitName() string
	ServiceName() string
	Config() *config.Config
	Logger() logger.Logger
	Clients() *Clients
	Storage() *postgres.Postgres
	MesServer() *rmqserver.Server
	MesClient() *rmqclient.Client
	Start(context.Context) error
	Stop(context.Context) error
}

type RepositoryImpl struct {
	unitName    string
	serviceName string
	config      *config.Config
	logger      logger.Logger
	clients     *Clients
	storage     *postgres.Postgres
	mesServ     *rmqserver.Server
	mesClient   *rmqclient.Client
	rmqRouter   map[string]rmqserver.CallHandler
	cancel      StackCancelFunc[CancelFunc]
}

func NewRepository(_ context.Context, rmqRouter map[string]rmqserver.CallHandler) (Repository, error) {
	var err error

	repo := &RepositoryImpl{}

	repo.unitName = "notify_web"
	repo.serviceName = "notify"

	repo.config = initTestConfig()

	repo.logger, err = logger.Setup(repo.config.Log.Level, nil)

	if err != nil {
		panic(err)
	}

	repo.rmqRouter = rmqRouter

	return repo, nil
}

func (r *RepositoryImpl) UnitName() string {
	return r.unitName
}

func (r *RepositoryImpl) ServiceName() string {
	return r.serviceName
}

func (r *RepositoryImpl) Config() *config.Config {
	return r.config
}

func (r *RepositoryImpl) Logger() logger.Logger {
	return r.logger
}

func (r *RepositoryImpl) Clients() *Clients {
	return r.clients
}

func (r *RepositoryImpl) Storage() *postgres.Postgres {
	return r.storage
}

func (r *RepositoryImpl) MesServer() *rmqserver.Server {
	return r.mesServ
}

func (r *RepositoryImpl) MesClient() *rmqclient.Client {
	return r.mesClient
}

func (r *RepositoryImpl) Start(ctx context.Context) (err error) {
	stackCancel := make(StackCancelFunc[CancelFunc], 0)

	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("recover repository", logger.AnyAttr("mes", rec))
			_ = stackCancel.Clear()
			err = ErrRepositoryImpl
		}
	}()

	cancelDB := initDatabaseIntegration(ctx)
	stackCancel.Push(func() { cancelDB() })
	cancelMes := initMessagingIntegration(ctx)
	stackCancel.Push(func() { cancelMes() })

	r.clients, err = NewClients(ctx, r.config)
	if err != nil {
		panic("error init grpc clients" + err.Error())
	}

	r.storage, err = postgres.New(ctx, r.config.Storage.URL)
	if err != nil {
		panic("storage connection error" + err.Error())
	}

	stackCancel.Push(func() { r.storage.Close() })

	err = r.storage.Ping(ctx)
	if err != nil {
		panic("storage ping error" + err.Error())
	}

	migrateUp(pathToMigrate, r.config.Storage.URL)

	r.mesClient, err = rmqclient.New(
		r.config.Messaging.URL,
		r.config.Messaging.Server.RPCExchange,
		r.config.Messaging.Client.RPCExchange,
		r.config.Messaging.Topics,
	)
	if err != nil {
		panic("messaging client connection error" + err.Error())
	}

	stackCancel.Push(func() { _ = r.mesClient.Shutdown() })

	r.mesServ, err = rmqserver.New(
		r.config.Messaging.URL,
		r.config.Messaging.Server.RPCExchange,
		r.config.Messaging.Topics,
		r.rmqRouter,
		r.logger,
	)
	if err != nil {
		panic("messaging server connection error" + err.Error())
	}

	stackCancel.Push(func() { _ = r.mesServ.Shutdown() })
	r.cancel = stackCancel

	return err
}

func (r *RepositoryImpl) Stop(_ context.Context) error {
	return r.cancel.Clear()
}

func migrateUp(migratePath, url string) {
	m, err := migrate.New(
		"file://"+migratePath,
		url,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		panic(err)
	}
}
