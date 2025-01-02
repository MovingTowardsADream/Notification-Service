package multi_handler

import (
	"context"
	"log/slog"
	"sync"
)

type MultiHandler struct {
	mu  *sync.Mutex
	out []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{out: handlers, mu: &sync.Mutex{}}
}

func (h *MultiHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *MultiHandler) Handle(
	ctx context.Context,
	r slog.Record, //nolint:gocritic
) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, destHandler := range h.out {
		err := destHandler.Handle(ctx, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	h2.out = make([]slog.Handler, len(h.out))
	for i, h := range h.out {
		h2.out[i] = h.WithGroup(name)
	}
	return &h2
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h
	h2.out = make([]slog.Handler, len(h.out))
	for i, h := range h.out {
		h2.out[i] = h.WithAttrs(attrs)
	}
	return &h2
}
