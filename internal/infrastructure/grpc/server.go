package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"Notification_Service/internal/infrastructure/grpc/notify"
	"Notification_Service/internal/infrastructure/grpc/users"
	"Notification_Service/internal/infrastructure/observ/metrics"
	"Notification_Service/internal/interfaces/middleware"
	"Notification_Service/pkg/logger"
	"Notification_Service/pkg/utils"
)

const (
	_defaultPort = 8080
)

type Server struct {
	gRPCServer *grpc.Server
	log        logger.Logger
	port       int
}

func New(log logger.Logger, m *metrics.Metrics, notifySender notify.SendersNotify, editInfo users.UserInfo, opts ...Option) *Server {
	s := &Server{
		log:  log,
		port: _defaultPort,
	}

	for _, opt := range opts {
		opt(s)
	}

	mw := middleware.New(log, m)

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.LoggingInterceptor,
			mw.RecoveryInterceptor,
			mw.MetricInterceptor,
		),
	)

	notify.Notify(gRPCServer, notifySender)
	users.Users(gRPCServer, editInfo)

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

	l, err := net.Listen("tcp", utils.FormatAddress("", s.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("gRPC server started", logger.AnyAttr("addr", l.Addr().String()))

	if err := s.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Shutdown() {
	const op = "gRPC - Server.Shutdown"

	s.log.Info("stopping gRPC server", logger.AnyAttr("port", s.port), logger.AnyAttr("op", op))

	s.gRPCServer.GracefulStop()
}
