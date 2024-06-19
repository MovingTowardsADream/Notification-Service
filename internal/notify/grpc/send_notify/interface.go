package grpc_send_notify

import (
	"Notification_Service/internal/entity"
	"context"
)

type NotifySend interface {
	SendNotifyForUser(ctx context.Context, notifyRequest *entity.RequestNotification) error
}
