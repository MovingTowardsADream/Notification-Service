package send_notify

import (
	"Notification_Service/internal/entity"
	"context"
)

type (
	UsersDataCommunication interface {
		GetUserCommunication(ctx context.Context, id string) (entity.UserCommunication, error)
	}

	NotifyGateway interface {
		CreateNotifyMessageOnRabbitMQ(ctx context.Context, notify entity.Notify) error
	}
)
