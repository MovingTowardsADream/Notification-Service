package middleware

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

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

	var grpcCode string

	if st, ok := status.FromError(err); ok {
		grpcCode = st.Code().String()
	}

	lm.Duration.With(prometheus.Labels{"method": info.FullMethod, "status": grpcCode}).Observe(time.Since(start).Seconds())

	return resp, err
}
