package main

import (
	"context"
	"log/slog"
)

type contextHandlerCtxKey struct{}

// AppendAttrsToContext returns a new context with one or more [slog.Attr]
// values embedded in it.
func AppendAttrsToContext(ctx context.Context, attrs ...slog.Attr) context.Context {
	existing := AttrsFromContext(ctx)
	return context.WithValue(ctx, contextHandlerCtxKey{}, append(existing, attrs...))
}

// AttrsFromContext returns [slog.Attr] that were embedded into the context.
func AttrsFromContext(ctx context.Context) []slog.Attr {
	existing, ok := ctx.Value(contextHandlerCtxKey{}).([]slog.Attr)
	if !ok {
		return nil
	}
	return existing
}

// ContextHandler is a context-aware [slog.Handler].
type ContextHandler struct {
	slog.Handler
}

// Enabled implements [slog.Handler].
func (h ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

// Handle implements [slog.Handler].
//
// Handle extracts [slog.Attr] values from ctx and adds it into the
// [slog.Record].
//
// See [AppendAttrsToContext] and [AttrsFromContext].
func (h ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(AttrsFromContext(ctx)...)
	return h.Handler.Handle(ctx, record)
}

// WithAttrs implements [slog.Handler].
func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

// WithGroup implements [slog.Handler].
func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{Handler: h.Handler.WithGroup(name)}
}
