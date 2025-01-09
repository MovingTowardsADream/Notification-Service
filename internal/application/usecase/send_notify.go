package usecase

import (
	"context"
	"errors"
	"fmt"

	repoErr "Notification_Service/internal/infrastructure/repository/errors"
	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/logger"
)

type UsersDataCommunication interface {
	GetUserCommunication(ctx context.Context, communication *dto.IdentificationUserCommunication) (*dto.UserCommunication, error)
}

type NotifyGateway interface {
	CreateMailNotify(ctx context.Context, notify *dto.MailDate) error
	CreatePhoneNotify(ctx context.Context, notify *dto.PhoneDate) error
}

type NotifySender struct {
	l             *logger.Logger
	usersDataComm UsersDataCommunication
	gateway       NotifyGateway
}

func NewNotifySender(l *logger.Logger, usersDataComm UsersDataCommunication, gateway NotifyGateway) *NotifySender {
	return &NotifySender{
		l:             l,
		usersDataComm: usersDataComm,
		gateway:       gateway,
	}
}

func (n *NotifySender) SendToUser(ctx context.Context, notifyRequest *dto.ReqNotification) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	userCommunication, err := n.usersDataComm.GetUserCommunication(ctxTimeout, &dto.IdentificationUserCommunication{ID: notifyRequest.UserID})

	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return ErrNotFound
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeout
		}

		n.l.Error("SendNotifyForUsers - n.usersDataComm.GetUserCommunication: ", n.l.Err(err))

		return err
	}

	if userCommunication.MailPref {
		mailNotify := &dto.MailDate{
			Mail:       userCommunication.Email,
			NotifyType: notifyRequest.NotifyType,
			Subject:    notifyRequest.Channels.Mail.Subject,
			Body:       notifyRequest.Channels.Mail.Body,
		}

		err = n.gateway.CreateMailNotify(ctxTimeout, mailNotify)

		if err != nil {
			return fmt.Errorf("SendNotifyForUser - n.gateway.CreateNotifyMessageOnRabbitMQ: %w", err)
		}
	}

	if userCommunication.PhonePref {
		phoneNotify := &dto.PhoneDate{
			Phone:      userCommunication.Phone,
			NotifyType: notifyRequest.NotifyType,
			Body:       notifyRequest.Channels.Phone.Body,
		}

		err = n.gateway.CreatePhoneNotify(ctxTimeout, phoneNotify)

		if err != nil {
			return fmt.Errorf("SendNotifyForUser - n.gateway.CreateNotifyMessageOnRabbitMQ: %w", err)
		}
	}

	return nil
}
