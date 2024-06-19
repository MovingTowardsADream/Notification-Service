package postgresdb

import (
	"Notification_Service/internal/entity"
	"Notification_Service/pkg/postgres"
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
)

type NotifyRepo struct {
	db *postgres.Postgres
}

func NewNotifyRepo(pg *postgres.Postgres) *NotifyRepo {
	return &NotifyRepo{pg}
}

func (r *NotifyRepo) GetUserCommunication(ctx context.Context, id string) (entity.UserCommunication, error) {
	var userCommunication entity.UserCommunication

	return userCommunication, nil
}

func (r *NotifyRepo) EditUserPreferences(ctx context.Context, preferences *notifyv1.UserPreferencesRequest) error {
	return nil
}
