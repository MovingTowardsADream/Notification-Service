package messaging

import (
	"context"
	"fmt"

	"Notification_Service/internal/interfaces/dto"
)

type NotifyGatewayMessaging interface {
	RemoteCall(ctx context.Context, handler string, request any) error
}

type NotifyGateway struct {
	mes NotifyGatewayMessaging
}

func NewNotifyGateway(mes NotifyGatewayMessaging) *NotifyGateway {
	return &NotifyGateway{mes}
}

func (gw *NotifyGateway) CreateMailNotify(ctx context.Context, notify *dto.MailDate) error {
	err := wrapper(ctx, func() error {
		return gw.mes.RemoteCall(ctx, "createNewMailNotify", notify)
	})

	if err != nil {
		return fmt.Errorf("NotifyGateway - CreateNotifyMailMessageOnRabbitMQ - gw.rmq.RemoteCall: %w", err)
	}

	return nil
}

func (gw *NotifyGateway) CreatePhoneNotify(ctx context.Context, notify *dto.PhoneDate) error {
	err := wrapper(ctx, func() error {
		return gw.mes.RemoteCall(ctx, "createNewPhoneNotify", notify)
	})

	if err != nil {
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
