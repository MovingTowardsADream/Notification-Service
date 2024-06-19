package edit_info

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
)

type UsersDataPreferences interface {
	EditUserPreferences(ctx context.Context, preferences *notifyv1.UserPreferencesRequest) error
}
