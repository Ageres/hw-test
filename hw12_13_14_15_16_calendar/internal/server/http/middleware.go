package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *AppServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
			s.logger.WithError(err).Error("get remote ip address", map[string]any{
				"ip": ip,
			})
		}

		userAgent := r.UserAgent()
		if userAgent == "" {
			userAgent = "-"
			s.logger.Warn("get user agent", map[string]any{"userAgent": "not found"})
		}

		rw := &responseWriter{w, http.StatusOK}

		requestId := utils.GenerateRequestID()
		ctx := context.WithValue(r.Context(), utils.RequestIDKey, requestId)

		logger := s.logger.With(map[string]any{"requestId": requestId})

		ctx = logger.SetLoggerToCtx(ctx)

		newR := r.WithContext(ctx)

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Header().Set("X-Request-ID", requestId)

		next.ServeHTTP(rw, newR)

		logger.Info("request", map[string]any{
			"ip":         ip,
			"method":     r.Method,
			"path":       r.URL.Path,
			"protocol":   r.Proto,
			"status":     rw.status,
			"latency_ms": time.Since(start).Milliseconds(),
			"user_agent": userAgent,
		})
	})
}
