package grpcserver

import (
	grpc_send_notify "Notification_Service/internal/grpc/send_notify"
	grpc_users "Notification_Service/internal/grpc/users"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

const (
	_defaultPort = ":8080"
)

type Server struct {
	log        *slog.Logger
	gRPCServer *grpc.Server

	port string
}

func New(log *slog.Logger, opts ...Option) *Server {

	s := &Server{
		log:  log,
		port: _defaultPort,
	}

	for _, opt := range opts {
		opt(s)
	}

	gRPCServer := grpc.NewServer()

	grpc_send_notify.SendNotify(gRPCServer)
	grpc_users.Users(gRPCServer)

	s.gRPCServer = gRPCServer

	return s
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	const op = "gRPC - Server.Run"

	l, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("gRPC server started", slog.String("addr", l.Addr().String()))

	if err := s.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Shutdown() {
	const op = "gRPC - Server.Shutdown"

	s.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.String("port", s.port))

	s.gRPCServer.GracefulStop()
}
