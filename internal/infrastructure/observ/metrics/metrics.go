package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Duration *prometheus.HistogramVec
}

func New(reg prometheus.Registerer, nameApp string) *Metrics {
	m := &Metrics{
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: nameApp,
			Name:      "req_duration_sec",
			Help:      "duration of the request",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"status", "method"}),
	}
	reg.MustRegister(m.Duration)

	return m
}
