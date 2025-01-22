package messaging

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"Notification_Service/internal/domain/models"
	"Notification_Service/internal/interfaces/convert"
	"Notification_Service/internal/interfaces/dto"
)

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
	tracer := otel.Tracer("NotifyGateway")
	ctx, span := tracer.Start(ctx, "CreateMailNotify")
	defer span.End()

	if mailNotify == nil {
		return fmt.Errorf("NotifyGateway - CreateMailNotify - notify is nil")
	}

	err := wrapper(ctx, func() error {
		return gw.mes.RemoteCall(ctx, "mail_notify", mailNotify.NotifyType, convert.MailDateToMailInfo(mailNotify))
	})

	if err != nil {
		return fmt.Errorf("NotifyGateway - CreateNotifyMailMessageOnRabbitMQ - gw.rmq.RemoteCall: %w", err)
	}

	return nil
}

func (gw *NotifyGateway) CreatePhoneNotify(ctx context.Context, phoneNotify *dto.PhoneDate) error {
	tracer := otel.Tracer("NotifyGateway")
	ctx, span := tracer.Start(ctx, "CreatePhoneNotify")
	defer span.End()

	if phoneNotify == nil {
		err := fmt.Errorf("NotifyGateway - CreateMailNotify - notify is nil")
		span.RecordError(err)
		return err
	}

	err := wrapper(ctx, func() error {
		return gw.mes.RemoteCall(ctx, "phone_notify", phoneNotify.NotifyType, convert.PhoneDateToPhoneInfo(phoneNotify))
	})

	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("NotifyGateway - CreateNotifyPhoneMessageOnRabbitMQ - gw.rmq.RemoteCall: %w", err)
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
