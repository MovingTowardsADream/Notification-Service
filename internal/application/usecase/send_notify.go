package usecase

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	repoerr "Notification_Service/internal/infrastructure/repository/errors"
	"Notification_Service/internal/interfaces/convert"
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
	tracer := otel.Tracer("NotifySender")
	ctx, span := tracer.Start(ctx, "SendToUser")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", notifyRequest.UserID))

	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	userCommunication, err := n.usersDataComm.GetUserCommunication(
		ctxTimeout,
		convert.ReqNotifyToIDUserCommunication(notifyRequest),
	)

	if err != nil {
		span.RecordError(err)

		if errors.Is(err, repoerr.ErrNotFound) {
			return ErrNotFound
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeout
		}

		n.l.Error(
			"SendNotifyForUsers - n.usersDataComm.GetUserCommunication",
			n.l.Err(err),
			logger.NewStrArgs("trace-id", ctx.Value("trace-id").(string)),
		)

		return err
	}

	if userCommunication.MailPref {
		mailNotify := convert.ToMailDate(notifyRequest, userCommunication)

		err = n.gateway.CreateMailNotify(ctxTimeout, mailNotify)

		if err != nil {
			span.RecordError(err)

			return fmt.Errorf("SendNotifyForUser - n.gateway.CreateNotifyMessageOnRabbitMQ: %w", err)
		}
	}

	if userCommunication.PhonePref {
		phoneNotify := convert.ToPhoneDate(notifyRequest, userCommunication)

		err = n.gateway.CreatePhoneNotify(ctxTimeout, phoneNotify)

		if err != nil {
			span.RecordError(err)

			return fmt.Errorf("SendNotifyForUser - n.gateway.CreateNotifyMessageOnRabbitMQ: %w", err)
		}
	}

	return nil
}
