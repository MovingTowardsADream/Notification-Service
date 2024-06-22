package grpcserver

import (
	"Notification_Service/internal/notify/grpc/send_notify"
	"Notification_Service/internal/notify/grpc/users"
	custom_validator "Notification_Service/pkg/validator"
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

func New(log *slog.Logger, notifySend grpc_send_notify.NotifySend, editInfo grpc_users.EditInfo, opts ...Option) *Server {

	s := &Server{
		log:  log,
		port: _defaultPort,
	}

	for _, opt := range opts {
		opt(s)
	}

	gRPCServer := grpc.NewServer()

	validator := custom_validator.NewCustomValidator()

	grpc_send_notify.SendNotify(gRPCServer, notifySend, validator)
	grpc_users.Users(gRPCServer, editInfo, validator)

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
