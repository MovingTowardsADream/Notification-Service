package usecase

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"Notification_Service/internal/interfaces/convert"
	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/logger"
)

type UsersDataCommunication interface {
	GetUserCommunication(ctx context.Context, communication *dto.IdentificationUserCommunication) (*dto.UserCommunication, error)
}

type NotifyProcessed interface {
	ProcessedNotify(ctx context.Context, processed *dto.ProcessedNotify) error
}

type NotifySender struct {
	l               logger.Logger
	usersDataComm   UsersDataCommunication
	notifyProcessed NotifyProcessed
}

func NewNotifySender(l logger.Logger, usersDataComm UsersDataCommunication, processed NotifyProcessed) *NotifySender {
	return &NotifySender{
		l:               l,
		usersDataComm:   usersDataComm,
		notifyProcessed: processed,
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
	span.SetAttributes(attribute.String("user.id", notifyRequest.RequestID))

	userCommunication, err := n.usersDataComm.GetUserCommunication(
		ctx,
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

	processedNotify := convert.ToProcessedNotify(notifyRequest, userCommunication)

	err = n.notifyProcessed.ProcessedNotify(ctx, processedNotify)

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

	return nil
}
