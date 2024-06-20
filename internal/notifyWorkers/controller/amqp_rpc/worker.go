package amqp_rpc

import (
	"Notification_Service/internal/entity"
	"Notification_Service/internal/notifyWorkers/usecase"
	rmq_server "Notification_Service/pkg/rabbitmq/server"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type notifyWorkerRoutes struct {
	w usecase.NotifyWorker
}

func newNotifyWorkerRoutes(routes map[string]rmq_server.CallHandler, w usecase.NotifyWorker) {
	r := &notifyWorkerRoutes{w}
	{
		routes["createNewMailNotify"] = r.createNewMailNotify()
		routes["createNewPhoneNotify"] = r.createNewPhoneNotify()
	}
}

func (r *notifyWorkerRoutes) createNewMailNotify() rmq_server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		var request entity.MailDate

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

func (r *notifyWorkerRoutes) createNewPhoneNotify() rmq_server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		var request entity.PhoneDate

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
