package middleware

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Notification_Service/pkg/logger"
)

type recoverMiddlewares struct {
	*logger.Logger
}

func (rm *recoverMiddlewares) RecoveryInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			rm.Error("panic occurred: ", slog.Any("panic", r))
			err = status.Errorf(codes.Internal, "panic occurred: %v", r)
		}
	}()

	return handler(ctx, req)
}
