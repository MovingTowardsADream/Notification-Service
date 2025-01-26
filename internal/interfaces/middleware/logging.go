package middleware

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"Notification_Service/pkg/logger"
)

type loggerMiddlewares struct {
	logger.Logger
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

	var grpcCode string

	if st, ok := status.FromError(err); ok {
		grpcCode = st.Code().String()
	}

	lm.Info("request end",
		logger.AnyAttr("method", info.FullMethod),
		logger.AnyAttr("status", grpcCode),
		logger.AnyAttr("trace-id", ctx.Value("trace-id").(string)),
		logger.AnyAttr("duration", time.Since(start)),
	)

	return resp, err
}

func generateTraceID() string {
	return uuid.New().String()
}
