package rmq_server

import (
	"Notification_Service/pkg/logger"
	rmq_rpc "Notification_Service/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log/slog"
	"runtime"
	"sync"
	"time"
)

const (
	_defaultWaitTime = 2 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second

	Success = "success"
)

var _defaultGoroutinesCount = runtime.NumCPU()

type CallHandler func(*amqp.Delivery) (interface{}, error)

type Server struct {
	serverExchange string
	conn           *rmq_rpc.Connection
	error          chan error
	stop           chan struct{}
	router         map[string]CallHandler

	timeout         time.Duration
	goroutinesCount int

	logger *slog.Logger

	rw       sync.RWMutex
	mistakes map[string]int
}

func New(url, serverExchange string, router map[string]CallHandler, l *slog.Logger, opts ...Option) (*Server, error) {
	cfg := rmq_rpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	s := &Server{
		serverExchange:  serverExchange,
		conn:            rmq_rpc.NewConnectionRabbitMQ(serverExchange, cfg),
		error:           make(chan error),
		stop:            make(chan struct{}),
		router:          router,
		timeout:         _defaultTimeout,
		goroutinesCount: _defaultGoroutinesCount,
		logger:          l,
	}

	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc server - NewServer - s.conn.AttemptConnect: %w", err)
	}

	return s, nil
}

func (s *Server) MustRun() {
	for i := 0; i < s.goroutinesCount; i++ {
		go s.consumer()
	}
}

func (s *Server) consumer() {
	for {
		select {
		case <-s.stop:
			return
		case d, opened := <-s.conn.Delivery:
			if !opened {
				s.reconnect()

				return
			}

			_ = d.Ack(false)

			s.serveCall(&d)
		}
	}
}

func (s *Server) republish(corrID, handler string, request interface{}) error {
	var (
		requestBody []byte
		err         error
	)

	if request != nil {
		requestBody, err = json.Marshal(request)
		if err != nil {
			return fmt.Errorf("publish - json.Marshal: %w", err)
		}
	}

	err = s.conn.Channel.Publish(s.serverExchange, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       s.serverExchange,
			Type:          handler,
			Body:          requestBody,
		})
	if err != nil {
		return fmt.Errorf("c.Channel.Publish: %w", err)
	}

	return nil
}

func (s *Server) serveCall(d *amqp.Delivery) {
	callHandler, ok := s.router[d.Type]
	if !ok {
		return
	}

	response, err := callHandler(d)

	if err != nil {
		count := s.getMistake(d.CorrelationId)

		if count >= 3 {
			s.deleteMistake(d.CorrelationId)
			return
		}

		s.addMistake(d.CorrelationId)

		err = s.republish(d.CorrelationId, d.Type, response)

		if err != nil {
			s.logger.Error("rmq_server - Server - serveCall - s.Republish", logger.Err(err))
			return
		}

		return
	}

	s.deleteMistake(d.CorrelationId)
}

func (s *Server) getMistake(corrID string) int {
	s.rw.RLock()
	count, ok := s.mistakes[corrID]
	s.rw.RUnlock()

	if !ok {
		return 0
	}
	return count
}

func (s *Server) addMistake(corrID string) {
	s.rw.Lock()
	s.mistakes[corrID]++
	s.rw.Unlock()
}

func (s *Server) deleteMistake(corrID string) {
	s.rw.Lock()
	delete(s.mistakes, corrID)
	s.rw.Unlock()
}

func (s *Server) reconnect() {
	close(s.stop)

	err := s.conn.AttemptConnect()
	if err != nil {
		s.error <- err
		close(s.error)

		return
	}

	s.stop = make(chan struct{})

	go s.consumer()
}

func (s *Server) Notify() <-chan error {
	return s.error
}

func (s *Server) Shutdown() error {
	select {
	case <-s.error:
		return nil
	default:
	}

	close(s.stop)
	time.Sleep(s.timeout)

	err := s.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("rmq_rpc server - Server - Shutdown - s.Connection.Close: %w", err)
	}

	return nil
}
