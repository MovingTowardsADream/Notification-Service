package grpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

// Logging
func (s *Server) loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Processing request
	resp, err := handler(ctx, req)

	// Записываем время окончания обработки запроса
	s.log.Info("Request end: ", info.FullMethod, "Request: ", fmt.Sprintf("%+v", req), " Duration: ", time.Since(start))

	return resp, err
}

// Recovery
func (s *Server) recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Error("panic occurred: ", slog.Any("panic", r))
			err = status.Errorf(codes.Internal, "panic occurred: %v", r)
		}
	}()

	return handler(ctx, req)
}
