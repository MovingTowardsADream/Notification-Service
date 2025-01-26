package gotests

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"google.golang.org/grpc"

	"Notification_Service/internal/infrastructure/grpc/notify"
	"Notification_Service/internal/infrastructure/grpc/users"
	"Notification_Service/internal/interfaces/middleware"
	"Notification_Service/pkg/logger"
)

const (
	_defaultMockServerPort = 8081
	_defaultMockServerHost = "127.0.0.1"
)

type api struct {
	notifySend notify.SendersNotify
	editInfo   users.EditInfo
}

type MockServer struct {
	port       int
	addr       string
	errCh      chan error
	cancelOnce *sync.Once
	log        *logger.Logger
	clients    api
	cancelFunc func()
}

type Option func(*MockServer)

func WithPort(port int) Option {
	return func(a *MockServer) {
		a.port = port
	}
}

func NewMockServer(log *logger.Logger, notifySend notify.SendersNotify, editInfo users.EditInfo, options ...Option) *MockServer {
	s := &MockServer{
		port:  _defaultMockServerPort,
		errCh: make(chan error),
		log:   log,
		clients: api{
			notifySend: notifySend,
			editInfo:   editInfo,
		},
	}

	for _, option := range options {
		option(s)
	}

	s.addr = fmt.Sprintf("%s:%d", _defaultMockServerHost, s.port)

	return s
}

func (s *MockServer) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("unable to listen on: %s, err: %w", s.addr, err)
	}
	s.log.Info("Starting MockServer on", slog.String("addr", listener.Addr().String()))
	return s.Serve(ctx, listener)
}

func (s *MockServer) Serve(ctx context.Context, listener net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)

	s.cancelOnce = &sync.Once{}
	s.cancelFunc = cancel

	go func(ctx context.Context, listener net.Listener) {
		<-ctx.Done()
		_ = listener.Close()
	}(ctx, listener)

	mw := middleware.New(s.log)

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.LoggingInterceptor,
			mw.RecoveryInterceptor,
		),
	)

	notify.Notify(gRPCServer, s.clients.notifySend)
	users.Users(gRPCServer, s.clients.editInfo)

	go func() {
		serv := gRPCServer.Serve(listener)
		if serv != nil {
			s.errCh <- serv
		}
		close(s.errCh)
	}()

	return nil
}

func (s *MockServer) Close() error {
	s.cancelOnce.Do(s.cancelFunc)
	return <-s.errCh
}

func (s *MockServer) Err() <-chan error {
	return s.errCh
}
