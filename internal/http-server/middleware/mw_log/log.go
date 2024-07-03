package mwlog

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logs for http-requests are enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			next.ServeHTTP(rw, r)

			entry.Info("request has been processed",
				slog.Int("status", rw.Status()),
				slog.String("size", fmt.Sprintf("%d bytes", rw.BytesWritten())),
				slog.String("duration", time.Since(startedAt).String()),
			)
		})
	}
}
