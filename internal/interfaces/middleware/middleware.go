package middleware

import (
	"Notification_Service/pkg/logger"
)

type Middlewares struct {
	recoverMiddlewares
	loggerMiddlewares
}

func New(l *logger.Logger) *Middlewares {
	return &Middlewares{
		recoverMiddlewares{l},
		loggerMiddlewares{l},
	}
}
