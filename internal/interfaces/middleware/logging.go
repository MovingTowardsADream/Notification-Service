package middleware

import (
	"context"
	"time"

	"github.com/google/uuid"
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

	ctx = context.WithValue(ctx, "trace-id", generateTraceID())

	resp, err := handler(ctx, req)

	respStatus := "ok"
	if err != nil {
		respStatus = "failed"
	}

	lm.Info("request end",
		logger.NewStrArgs("method", info.FullMethod),
		logger.NewStrArgs("status", respStatus),
		logger.NewStrArgs("trace-id", ctx.Value("trace-id").(string)),
		logger.NewDurationArgs("duration", time.Since(start)),
	)

	return resp, err
}

func generateTraceID() string {
	return uuid.New().String()
}
