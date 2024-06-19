package grpc_send_notify

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
)

type NotifySend interface {
	SendNotifyForUser(ctx context.Context, notifyRequest *notifyv1.SendMessageRequest) error
}
