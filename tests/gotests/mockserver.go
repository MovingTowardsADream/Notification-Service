package gotests

import (
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"

	"Notification_Service/internal/infrastructure/grpc/notify"
	"Notification_Service/internal/infrastructure/grpc/users"
	"Notification_Service/pkg/utils"
)

const (
	_defaultMockServerPort = 8082
	_defaultMockServerHost = ""
)

type MockServer struct {
	port       int
	addr       string
	errCh      chan error
	cancelOnce *sync.Once
	gRPCServer *grpc.Server
	cancelFunc func()
}

type Option func(*MockServer)

func WithPort(port int) Option {
	return func(a *MockServer) {
		a.port = port
	}
}

func NewMockServer(notifySend notify.SendersNotify, editInfo users.EditInfo, options ...Option) *MockServer {
	s := &MockServer{
		port:  _defaultMockServerPort,
		errCh: make(chan error),
	}
	for _, option := range options {
		option(s)
	}

	s.addr = utils.FormatAddress(_defaultMockServerHost, s.port)

	gRPCServer := grpc.NewServer()

	notify.Notify(gRPCServer, notifySend)
	users.Users(gRPCServer, editInfo)

	s.gRPCServer = gRPCServer

	return s
}

func (s *MockServer) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("unable to listen on: %s, err: %w", s.addr, err)
	}

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

	go func() {
		serv := s.gRPCServer.Serve(listener)

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
