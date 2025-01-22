package trace

import (
	"time"
)

type Option func(provider *TracesProvider)

func Enabled(enabled bool) Option {
	return func(provider *TracesProvider) {
		provider.enabled = enabled
	}
}

func InitialInterval(interval time.Duration) Option {
	return func(provider *TracesProvider) {
		provider.initialInterval = interval
	}
}

func MaxInterval(interval time.Duration) Option {
	return func(provider *TracesProvider) {
		provider.maxInterval = interval
	}
}

func MaxElapsedTime(elapsedTime time.Duration) Option {
	return func(provider *TracesProvider) {
		provider.maxElapsedTime = elapsedTime
	}
}
