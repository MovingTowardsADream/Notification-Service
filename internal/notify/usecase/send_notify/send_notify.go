package send_notify

import (
	"Notification_Service/internal/entity"
	"Notification_Service/internal/notify/repository/repository_erros"
	"Notification_Service/internal/notify/usecase/usecase_errors"
	"Notification_Service/pkg/logger"
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

		n.l.Error("SendNotifyForUsers - n.usersDataComm.GetUserCommunication: ", logger.Err(err))

		return err
	}

	var mail *entity.MailDate
	var phone *entity.PhoneDate

	if userCommunication.MailPref {
		mail = &entity.MailDate{
			Mail:       userCommunication.Email,
			NotifyType: notifyRequest.NotifyType,
			Subject:    notifyRequest.Channels.Mail.Subject,
			Body:       notifyRequest.Channels.Mail.Body,
		}
	}

	if userCommunication.PhonePref {
		phone = &entity.PhoneDate{
			Phone:      userCommunication.Phone,
			NotifyType: notifyRequest.NotifyType,
			Body:       notifyRequest.Channels.Phone.Body,
		}
	}

	notify := entity.Notify{
		MailDate:  mail,
		PhoneDate: phone,
	}

	// TODO Gateway

	fmt.Println(notify.PhoneDate, notify.MailDate)

	return nil
}
