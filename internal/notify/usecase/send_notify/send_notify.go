package send_notify

import (
	"Notification_Service/internal/entity"
	"Notification_Service/internal/notify/repository/repository_erros"
	"Notification_Service/internal/notify/usecase/usecase_errors"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

const (
	_defaultTimeout = 5 * time.Second
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
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	userCommunication, err := n.usersDataComm.GetUserCommunication(ctxTimeout, notifyRequest.UserId)

	if err != nil {
		if errors.Is(err, usecase_errors.ErrTimeout) {
			return usecase_errors.ErrTimeout
		} else if errors.Is(err, repository_erros.ErrNotFound) {
			return usecase_errors.ErrNotFound
		}

		// TODO logging error

		return fmt.Errorf("UseCase - NotifySend - SendNotifyForUsers - n.usersDataComm.GetUserCommunication: %w", err)
	}

	// TODO Gateway

	_ = userCommunication

	return nil
}
