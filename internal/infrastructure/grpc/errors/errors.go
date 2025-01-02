package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDeadlineExceeded = status.Error(codes.DeadlineExceeded, "deadline exceeded")
	ErrInternalServer   = status.Error(codes.Internal, "internal server error")
	ErrNotFound         = status.Error(codes.NotFound, "object not found")
)
