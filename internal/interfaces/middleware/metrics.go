package middleware

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"Notification_Service/internal/infrastructure/observ/metrics"
)

type metricMiddlewares struct {
	*metrics.Metrics
}

func (lm *metricMiddlewares) MetricInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	respStatus := "ok"
	if err != nil {
		respStatus = "failed"
	}

	lm.Duration.With(prometheus.Labels{"method": info.FullMethod, "status": respStatus}).Observe(time.Since(start).Seconds())

	return resp, err
}
