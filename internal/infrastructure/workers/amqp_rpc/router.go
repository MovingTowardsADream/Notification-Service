package amqp_rpc

import (
	"Notification_Service/internal/infrastructure/smtp"

	rmqServer "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
)

func NewRouter(r smtp.NotifyWorker) map[string]rmqServer.CallHandler {
	routes := make(map[string]rmqServer.CallHandler)
	{
		newNotifyWorkerRoutes(routes, r)
	}

	return routes
}
