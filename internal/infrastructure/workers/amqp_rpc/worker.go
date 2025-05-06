package amqprpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/interfaces/dto"
)

//go:generate mockgen -source=worker.go -destination=../../../../tests/gotests/mocks/sender_mocks.go -package=mocks

type SenderMail interface {
	SendMailLetter(ctx context.Context, notify dto.MailInfo) error
}

type SenderPhone interface {
	SendPhoneSMS(ctx context.Context, notify dto.PhoneInfo) error
}

type notifyWorkerRoutes struct {
	sm SenderMail
	sp SenderPhone
}

func newNotifyWorkerRoutes(routes map[string]rmqserver.CallHandler, sm SenderMail, sp SenderPhone) {
	r := &notifyWorkerRoutes{sm, sp}
	//nolint:gocritic // increased readability
	{
		routes["notify.mail"] = r.createNewMailNotify()
		routes["notify.phone"] = r.createNewPhoneNotify()
	}
}

func (r *notifyWorkerRoutes) createNewMailNotify() rmqserver.CallHandler {
	const op = "amqp_rpc - createNewMailNotify"

	return func(d *amqp.Delivery) (any, error) {
		var request dto.MailInfo

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("%s - json.Unmarshal: %w", op, err)
		}

		err := r.sm.SendMailLetter(context.Background(), request)
		if err != nil {
			return request, fmt.Errorf("%s - SendMailLetter: %w", op, err)
		}

		return nil, nil
	}
}

func (r *notifyWorkerRoutes) createNewPhoneNotify() rmqserver.CallHandler {
	const op = "amqp_rpc - createNewPhoneNotify"

	return func(d *amqp.Delivery) (any, error) {
		var request dto.PhoneInfo

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("%s - json.Unmarshal: %w", op, err)
		}

		err := r.sp.SendPhoneSMS(context.Background(), request)
		if err != nil {
			return request, fmt.Errorf("%s - SendPhoneSMS: %w", op, err)
		}

		return nil, nil
	}
}
