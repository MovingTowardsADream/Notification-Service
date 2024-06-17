package grpcserver

import "time"

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
