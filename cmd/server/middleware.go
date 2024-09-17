package main

import (
	"log/slog"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
)

type Middleware = func(next http.Handler, meta routeMetadata) http.Handler

func RequestIDMiddleware() Middleware {
	return func(next http.Handler, meta routeMetadata) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("request-id")
			if reqID == "" {
				reqID = uuid.NewString()
			}

			w.Header().Set("request-id", reqID)

			ctx := r.Context()
			ctx = AppendAttrsToContext(ctx, slog.String("request_id", reqID))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoggingMiddleware(logger *slog.Logger) Middleware {
	logger = logger.With("event_type", "middleware.logging")
	return func(next http.Handler, meta routeMetadata) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			metrics := httpsnoop.CaptureMetrics(next, w, r)

			log := logger.InfoContext
			if metrics.Code >= 400 {
				log = logger.WarnContext
			}

			log(
				r.Context(), "request",
				"route", meta.Name,
				"code", metrics.Code,
				"duration", metrics.Duration.String(),
				"written_bytes", metrics.Written,
			)
		})
	}
}
