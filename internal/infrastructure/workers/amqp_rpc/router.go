package amqprpc

import (
	"Notification_Service/internal/infrastructure/smtp"

	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
)

func NewRouter(r smtp.NotifyWorker) map[string]rmqserver.CallHandler {
	routes := make(map[string]rmqserver.CallHandler)
	{
		newNotifyWorkerRoutes(routes, r)
	}

	return routes
}
