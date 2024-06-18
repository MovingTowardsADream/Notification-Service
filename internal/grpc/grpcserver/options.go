package grpcserver

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}
