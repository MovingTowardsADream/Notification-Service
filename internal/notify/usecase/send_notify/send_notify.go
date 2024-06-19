package send_notify

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"log/slog"
)

type NotifySend struct {
	l             *slog.Logger
	usersDataComm UsersDataCommunication
}

func New(l *slog.Logger, usersDataComm UsersDataCommunication) *NotifySend {
	return &NotifySend{
		l:             l,
		usersDataComm: usersDataComm,
	}
}

func (n *NotifySend) SendNotifyForUser(ctx context.Context, notifyRequest *notifyv1.SendMessageRequest) error {

	userCommunication, err := n.usersDataComm.GetUserCommunication(ctx, notifyRequest.UserId)

	_ = userCommunication

	return err
}
