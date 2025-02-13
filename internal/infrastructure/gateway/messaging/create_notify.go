package messaging

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"Notification_Service/internal/domain/models"
	"Notification_Service/internal/interfaces/convert"
	"Notification_Service/internal/interfaces/dto"
)

const traceName = "NotifyGateway"

type NotifyGatewayMessaging interface {
	RemoteCall(ctx context.Context, handler string, priority models.NotifyType, request any) error
}

type NotifyGateway struct {
	mes NotifyGatewayMessaging
}

func NewNotifyGateway(mes NotifyGatewayMessaging) *NotifyGateway {
	return &NotifyGateway{mes}
}

func (gw *NotifyGateway) CreateMailNotify(ctx context.Context, mailNotify *dto.MailDate) error {
	const op = "NotifyGateway - CreateMailNotify"
	const spanName = "CreateMailNotify"

	tracer := otel.Tracer(traceName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	if mailNotify == nil {
		return fmt.Errorf("%s - notify is nil", op)
	}

	err := wrapper(ctx, func() error {
		return gw.mes.RemoteCall(ctx, "mail_notify", mailNotify.NotifyType, convert.MailDateToMailInfo(mailNotify))
	})

	if err != nil {
		return fmt.Errorf("%s - gw.rmq.RemoteCall: %w", op, err)
	}

	return nil
}

func (gw *NotifyGateway) CreatePhoneNotify(ctx context.Context, phoneNotify *dto.PhoneDate) error {
	const op = "NotifyGateway - CreatePhoneNotify"
	const spanName = "CreatePhoneNotify"

	tracer := otel.Tracer(traceName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	if phoneNotify == nil {
		err := fmt.Errorf("%s - notify is nil", op)
		span.RecordError(err)
		return err
	}

	err := wrapper(ctx, func() error {
		return gw.mes.RemoteCall(ctx, "phone_notify", phoneNotify.NotifyType, convert.PhoneDateToPhoneInfo(phoneNotify))
	})

	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("%s - gw.rmq.RemoteCall: %w", op, err)
	}

	return nil
}

func wrapper(ctx context.Context, f func() error) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
