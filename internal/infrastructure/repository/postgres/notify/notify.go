package notify

import (
	"context"

	"Notification_Service/internal/infrastructure/repository/postgres"
)

const (
	usersTable  = "users"
	notifyTable = "notifications"
)

const tracerName = "notifyRepo"

type NotifyRepo struct {
	storage *postgres.Postgres
}

func NewNotifyRepo(storage *postgres.Postgres) *NotifyRepo {
	return &NotifyRepo{storage: storage}
}

func (r *NotifyRepo) ProcessedNotify(ctx context.Context, processed &dto.) error {

}
