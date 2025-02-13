package rmqserver

import "time"

type Option func(*Server)

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func GoroutinesCount(count int) Option {
	return func(s *Server) {
		s.goroutinesCount = count
	}
}

func ConnWaitTime(timeout time.Duration) Option {
	return func(s *Server) {
		s.conn.WaitTime = timeout
	}
}

func ConnAttempts(attempts int) Option {
	return func(s *Server) {
		s.conn.Attempts = attempts
	}
}
