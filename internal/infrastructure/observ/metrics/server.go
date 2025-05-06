package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"Notification_Service/pkg/logger"
)

const (
	defaultReadHeaderTimeout = 30 * time.Second
	defaultReadTimeout       = 60 * time.Second
	defaultWriteTimeout      = 60 * time.Second
	defaultIdleTimeout       = 120 * time.Second
)

type Server struct {
	httpServer *http.Server
	l          logger.Logger
}

func NewMetricsServer(address string, reg *prometheus.Registry, l logger.Logger) *Server {
	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:              address,
			Handler:           pMux,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ReadTimeout:       defaultReadTimeout,
			WriteTimeout:      defaultWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
		},
		l: l,
	}
}

func (s *Server) Run() error {
	s.l.Info("starting metrics server", logger.AnyAttr("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
