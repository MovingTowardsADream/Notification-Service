package amqprpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	"Notification_Service/internal/infrastructure/smtp"

	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
	"Notification_Service/internal/interfaces/dto"
)

type notifyWorkerRoutes struct {
	w smtp.NotifyWorker
}

func newNotifyWorkerRoutes(routes map[string]rmqserver.CallHandler, w smtp.NotifyWorker) {
	r := &notifyWorkerRoutes{w}
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

		err := r.w.CreateNewMailNotify(context.Background(), request)
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

		err := r.w.CreateNewPhoneNotify(context.Background(), request)
		if err != nil {
			return request, fmt.Errorf("amqp_rpc - notifyWorkerRoutes - createNewPhoneNotify - r.w.CreateNewNotify: %w", err)
		}

		return nil, nil
	}
}
