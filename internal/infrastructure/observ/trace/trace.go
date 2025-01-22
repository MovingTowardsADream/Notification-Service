package trace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

const (
	_defaultTraceConfigEnabled   = true
	_defaultTraceInitialInterval = 500 * time.Millisecond
	_defaultTraceMaxInterval     = 5 * time.Second
	_defaultTraceMaxElapsedTime  = 1 * time.Minute
)

type TracesProvider struct {
	Exporter *sdktrace.TracerProvider

	enabled         bool
	initialInterval time.Duration
	maxInterval     time.Duration
	maxElapsedTime  time.Duration
}

func New(ctx context.Context, host, appName string, opts ...Option) (*TracesProvider, error) {
	provider := &TracesProvider{
		enabled:         _defaultTraceConfigEnabled,
		initialInterval: _defaultTraceInitialInterval,
		maxInterval:     _defaultTraceMaxInterval,
		maxElapsedTime:  _defaultTraceMaxElapsedTime,
	}

	for _, opt := range opts {
		opt(provider)
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(host),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         provider.enabled,
			InitialInterval: provider.initialInterval,
			MaxInterval:     provider.maxInterval,
			MaxElapsedTime:  provider.maxElapsedTime,
		}),
	)
	if err != nil {
		return nil, ErrInitTraceExporter
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		)),
	)

	otel.SetTracerProvider(tp)

	provider.Exporter = tp

	return provider, nil
}

func (t *TracesProvider) Close(ctx context.Context) error {
	return t.Exporter.Shutdown(ctx)
}
