package send_notify

import (
	"Notification_Service/internal/entity"
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

func (n *NotifySend) SendNotifyForUser(ctx context.Context, notifyRequest *entity.RequestNotification) error {

	userCommunication, err := n.usersDataComm.GetUserCommunication(ctx, notifyRequest.UserId)

	_ = userCommunication

	return err
}
