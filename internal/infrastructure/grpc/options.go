package grpc

type Option func(*Server)

func Port(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}
