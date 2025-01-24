package middleware

import (
	"Notification_Service/internal/infrastructure/observ/metrics"
	"Notification_Service/pkg/logger"
)

type Middlewares struct {
	recoverMiddlewares
	loggerMiddlewares
	metricMiddlewares
}

func New(l *logger.Logger, m *metrics.Metrics) *Middlewares {
	return &Middlewares{
		recoverMiddlewares{l},
		loggerMiddlewares{l},
		metricMiddlewares{m},
	}
}
