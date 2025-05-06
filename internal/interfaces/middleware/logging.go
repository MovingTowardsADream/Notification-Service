package middleware

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"Notification_Service/internal/interfaces/dto"
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
	var (
		traceIDKey = dto.ContextKey("trace-id")
		grpcCode   string
	)

	start := time.Now()

	ctxWithValue := context.WithValue(ctx, traceIDKey, generateTraceID())

	resp, err := handler(ctxWithValue, req)

	if st, ok := status.FromError(err); ok {
		grpcCode = st.Code().String()
	}

	lm.Info("request end",
		logger.AnyAttr("method", info.FullMethod),
		logger.AnyAttr("status", grpcCode),
		logger.AnyAttr("trace-id", ctxWithValue.Value("trace-id").(string)),
		logger.AnyAttr("duration", time.Since(start)),
	)

	return resp, err
}

func generateTraceID() string {
	return uuid.New().String()
}
