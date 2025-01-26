package rmqserver

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/streadway/amqp"

	rmqrpc "Notification_Service/internal/infrastructure/messaging/rabbitmq"
	"Notification_Service/pkg/logger"
)

const (
	_defaultWaitTime = 2 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

var _defaultGoroutinesCount = runtime.NumCPU()

type CallHandler func(*amqp.Delivery) (any, error)

type Server struct {
	serverExchange string
	conn           *rmqrpc.Connection
	error          chan error
	stop           chan struct{}
	router         map[string]CallHandler

	timeout         time.Duration
	goroutinesCount int

	logger logger.Logger

	rw       sync.RWMutex
	mistakes map[string]int

	once sync.Once
}

func New(url, serverExchange string, topics []string, router map[string]CallHandler, l logger.Logger, opts ...Option) (*Server, error) {
	cfg := rmqrpc.Params{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	s := &Server{
		serverExchange:  serverExchange,
		conn:            rmqrpc.NewConnection(serverExchange, topics, cfg),
		error:           make(chan error),
		stop:            make(chan struct{}),
		router:          router,
		timeout:         _defaultTimeout,
		goroutinesCount: _defaultGoroutinesCount,
		logger:          l,
		mistakes:        make(map[string]int),
	}

	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnect(s.conn.ConnectReader())
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

func (s *Server) topicConsume(deliveryChan <-chan amqp.Delivery) {
	for {
		select {
		case <-s.stop:
			return
		case d, opened := <-deliveryChan:
			if !opened {
				s.once.Do(
					func() {
						go func() {
							s.logger.Warn("channel for topic closed. reconnecting...", logger.AnyAttr("topic", d.Type))
							s.reconnect()
						}()
					},
				)
				return
			}

			_ = d.Ack(false)

			s.serveCall(&d)
		}
	}
}

func (s *Server) consumer() {
	for _, deliveryChan := range s.conn.Deliveries {
		go s.topicConsume(deliveryChan)
	}
}

func (s *Server) republish(corrID, handler string, priority uint8, request any) error {
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

	err = s.conn.Channel.Publish(handler, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       s.serverExchange,
			Type:          handler,
			Body:          requestBody,
			Priority:      priority,
		})
	if err != nil {
		return fmt.Errorf("c.Channel.Publish: %w", err)
	}

	return nil
}

func (s *Server) serveCall(d *amqp.Delivery) {
	callHandler, ok := s.router[d.Type]
	if !ok {
		s.logger.Error("no handlers found for topic", logger.AnyAttr("topic", d.Type))
		return
	}
	response, err := callHandler(d)
	if err != nil {
		s.addMistake(d.CorrelationId)
		count := s.getMistake(d.CorrelationId)
		if count >= 3 {
			s.deleteMistake(d.CorrelationId)
			s.logger.Error(
				"max retry limit reached for message",
				logger.AnyAttr("id", d.CorrelationId),
				logger.AnyAttr("topic", d.Type),
				logger.AnyAttr("mistakes", count),
			)

			return
		}

		err = s.republish(d.CorrelationId, d.Type, d.Priority, response)
		if err != nil {
			s.logger.Error("rmq_server - Server - serveCall - s.republish", s.logger.Err(err))
			return
		}

		return
	}

	s.deleteMistake(d.CorrelationId)
}

func (s *Server) reconnect() {
	close(s.stop)

	err := s.conn.AttemptConnect(s.conn.ConnectReader())
	if err != nil {
		s.error <- err
		close(s.error)

		return
	}

	s.stop = make(chan struct{})
	s.once = sync.Once{}

	go s.consumer()
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
