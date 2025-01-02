package amqp_rpc

import (
	"Notification_Service/internal/infrastructure/smtp"

	rmq_server "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
)

func NewRouter(r smtp.NotifyWorker) map[string]rmq_server.CallHandler {
	routes := make(map[string]rmq_server.CallHandler)
	{
		newNotifyWorkerRoutes(routes, r)
	}

	return routes
}
