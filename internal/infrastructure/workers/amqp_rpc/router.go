package amqprpc

import (
	rmqserver "Notification_Service/internal/infrastructure/messaging/rabbitmq/server"
)

func NewRouter(sm SenderMail, sp SenderPhone) map[string]rmqserver.CallHandler {
	routes := make(map[string]rmqserver.CallHandler)
	{
		newNotifyWorkerRoutes(routes, sm, sp)
	}

	return routes
}
