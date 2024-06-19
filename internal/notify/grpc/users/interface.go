package grpc_users

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
)

type EditInfo interface {
	EditUserPreferences(ctx context.Context, preferences *notifyv1.UserPreferencesRequest) error
}
