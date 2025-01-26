package usecase

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

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
	l             logger.Logger
	usersDataComm UsersDataCommunication
	gateway       NotifyGateway
}

func NewNotifySender(l logger.Logger, usersDataComm UsersDataCommunication, gateway NotifyGateway) *NotifySender {
	return &NotifySender{
		l:             l,
		usersDataComm: usersDataComm,
		gateway:       gateway,
	}
}

func (n *NotifySender) SendToUser(ctx context.Context, notifyRequest *dto.ReqNotification) error {
	const op = "SendNotifyForUsers"

	const (
		tracerName = "NotifySender"
		spanName   = "SendToUser"
	)

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
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

		if ok, err := mappingErrors(err); ok {
			return err
		}

		n.l.Error(
			op,
			n.l.Err(err),
			logger.AnyAttr("trace-id", ctx.Value("trace-id").(string)),
		)

		return err
	}

	if userCommunication.MailPref {
		mailNotify := convert.ToMailDate(notifyRequest, userCommunication)

		err = n.gateway.CreateMailNotify(ctxTimeout, mailNotify)

		if err != nil {
			span.RecordError(err)

			return fmt.Errorf("%s - n.gateway.CreateNotifyMessageOnRabbitMQ: %w", op, err)
		}
	}

	if userCommunication.PhonePref {
		phoneNotify := convert.ToPhoneDate(notifyRequest, userCommunication)

		err = n.gateway.CreatePhoneNotify(ctxTimeout, phoneNotify)

		if err != nil {
			span.RecordError(err)

			return fmt.Errorf("%s - n.gateway.CreateNotifyMessageOnRabbitMQ: %w", op, err)
		}
	}

	return nil
}
