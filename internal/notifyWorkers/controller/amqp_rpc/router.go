package amqp_rpc

import (
	"Notification_Service/internal/notifyWorkers/usecase"
	rmq_server "Notification_Service/pkg/rabbitmq/server"
)

func NewRouter(r usecase.NotifyWorker) map[string]rmq_server.CallHandler {
	routes := make(map[string]rmq_server.CallHandler)
	{
		newNotifyWorkerRoutes(routes, r)
	}

	return routes
}
