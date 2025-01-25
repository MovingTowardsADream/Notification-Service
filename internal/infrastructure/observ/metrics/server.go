package metrics

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"Notification_Service/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	l          *logger.Logger
}

func NewMetricsServer(address string, reg *prometheus.Registry, l *logger.Logger) *Server {
	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: pMux,
		},
		l: l,
	}
}

func (s *Server) Run() error {
	s.l.Info("starting metrics server", slog.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
