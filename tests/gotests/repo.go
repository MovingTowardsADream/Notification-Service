package gotests

import (
	"context"

	"Notification_Service/internal/infrastructure/config"
	"Notification_Service/internal/infrastructure/repository/postgres"
)

type Repository interface {
	UnitName() string
	ServiceName() string
	Config() *config.Config
	Clients() *Clients
	Postgres() *postgres.Postgres
	Start(context.Context) error
	Stop(context.Context) error
}

type RepositoryImpl struct {
	unitName    string
	serviceName string
	config      *config.Config
	clients     *Clients
	storage     *postgres.Postgres
}

func NewRepository() (Repository, error) {
	repo := &RepositoryImpl{}
	repo.unitName = "notify_web"
	repo.serviceName = "notify"

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

func (r *RepositoryImpl) Postgres() *postgres.Postgres {
	return r.storage
}

func (r *RepositoryImpl) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
