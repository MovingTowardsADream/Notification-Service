package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
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

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	const op = "grpcserver.Run"

	l, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := s.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Shutdown() {
	const op = "grpcserver.Stop"

	s.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.String("port", s.port))

	s.gRPCServer.GracefulStop()
}
