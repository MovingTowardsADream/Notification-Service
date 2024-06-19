package grpc_users

import (
	"Notification_Service/internal/entity"
	"context"
)

type EditInfo interface {
	EditUserPreferences(ctx context.Context, preferences *entity.UserPreferences) error
}
