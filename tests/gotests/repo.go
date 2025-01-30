package gotests

import (
	"context"

	"Notification_Service/internal/infrastructure/config"
	rmqclient "Notification_Service/internal/infrastructure/messaging/rabbitmq/client"
	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/pkg/logger"
)

const (
	_defaultConfigPath = "./configs/testing.yaml"
	_defaultEnvPath    = ".env"
)

type Repository interface {
	UnitName() string
	ServiceName() string
	Config() *config.Config
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
	clients     *Clients
	storage     *postgres.Postgres
	mesServ     *rmqserver.Server
	mesClient   *rmqclient.Client
}

func NewRepository(ctx context.Context, rmqRouter map[string]rmqserver.CallHandler) (Repository, error) {
	var err error

	repo := &RepositoryImpl{}

	repo.unitName = "notify_web"
	repo.serviceName = "notify"

	repo.config = config.MustLoadPath(_defaultConfigPath, _defaultEnvPath)

	log, err := logger.Setup(_defaultEnvPath, nil)

	if err != nil {
		panic(err)
	}

	repo.clients, err = NewClients(ctx, repo.config)

	if err != nil {
		panic(err)
	}

	runDatabaseIntegration(ctx, repo.config)
	repo.storage, err = postgres.New(ctx, repo.config.Storage.URL)

	if err != nil {
		panic("storage connection error" + err.Error())
	}

	runMessagingIntegration(ctx, repo.config)

	repo.mesClient, err = rmqclient.New(
		repo.config.Messaging.URL,
		repo.config.Messaging.Server.RPCExchange,
		repo.config.Messaging.Client.RPCExchange,
		repo.config.Messaging.Topics,
	)

	if err != nil {
		panic("messaging client connection error" + err.Error())
	}

	repo.mesServ, err = rmqserver.New(
		repo.config.Messaging.URL,
		repo.config.Messaging.Server.RPCExchange,
		repo.config.Messaging.Topics,
		rmqRouter,
		log,
	)

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

func (r *RepositoryImpl) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
