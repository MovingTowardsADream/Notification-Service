package edit_info

import (
	"Notification_Service/internal/entity"
	"context"
)

type UsersDataPreferences interface {
	EditUserPreferences(ctx context.Context, preferences *entity.UserPreferences) error
}
