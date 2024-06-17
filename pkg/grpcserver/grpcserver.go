package grpcserver

import (
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

const (
	_defaultPort            = ":8080"
	_defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	log        *slog.Logger
	gRPCServer *grpc.Server

	shutdownTimeout time.Duration
	port            string
}

func New(log *slog.Logger, opts ...Option) *Server {

	s := &Server{
		log:             log,
		shutdownTimeout: _defaultShutdownTimeout,
		port:            _defaultPort,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
