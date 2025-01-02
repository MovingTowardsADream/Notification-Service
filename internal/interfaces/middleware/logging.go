package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"Notification_Service/pkg/logger"
)

type loggerMiddlewares struct {
	*logger.Logger
}

func (lm *loggerMiddlewares) LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	respStatus := "ok"
	if err != nil {
		respStatus = "failed"
	}

	lm.Info("Request end: ", info.FullMethod, "Status: ", respStatus, "Duration: ", time.Since(start))

	return resp, err
}
