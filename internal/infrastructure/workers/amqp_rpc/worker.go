package amqprpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/interfaces/dto"
)

type SenderMail interface {
	SendMailLetter(ctx context.Context, notify dto.MailDate) error
}

type SenderPhone interface {
	SendPhoneSMS(ctx context.Context, notify dto.PhoneDate) error
}

type notifyWorkerRoutes struct {
	sm SenderMail
	sp SenderPhone
}

func newNotifyWorkerRoutes(routes map[string]rmqserver.CallHandler, sm SenderMail, sp SenderPhone) {
	r := &notifyWorkerRoutes{sm, sp}
	{
		routes["mail_notify"] = r.createNewMailNotify()
		routes["phone_notify"] = r.createNewPhoneNotify()
	}
}

func (r *notifyWorkerRoutes) createNewMailNotify() rmqserver.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		var request dto.MailDate

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("amqp_rpc - notifyWorkerRoutes - createNewMailNotify - json.Unmarshal: %w", err)
		}

		err := r.sm.SendMailLetter(context.Background(), request)
		if err != nil {
			return request, fmt.Errorf("amqp_rpc - notifyWorkerRoutes - createNewMailNotify - r.w.CreateNewNotify: %w", err)
		}

		return nil, nil
	}
}

func (r *notifyWorkerRoutes) createNewPhoneNotify() rmqserver.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		var request dto.PhoneDate

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("amqp_rpc - notifyWorkerRoutes - createNewPhoneNotify - json.Unmarshal: %w", err)
		}

		err := r.sp.SendPhoneSMS(context.Background(), request)
		if err != nil {
			return request, fmt.Errorf("amqp_rpc - notifyWorkerRoutes - createNewPhoneNotify - r.w.CreateNewNotify: %w", err)
		}

		return nil, nil
	}
}
