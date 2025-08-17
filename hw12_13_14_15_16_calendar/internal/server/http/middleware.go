package internalhttp

import (
	"net"
	"net/http"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
)

type MyMiddleware struct {
}

func (s *AppServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()
		ctx = lg.SetLogger(ctx, s.logger)
		r.WithContext(ctx)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			s.logger.Error("get remote ip address", map[string]any{"error": err})
			ip = r.RemoteAddr
		}

		userAgent := "-"
		if len(r.UserAgent()) > 0 {
			userAgent = r.UserAgent()
		} else {
			s.logger.Warn("get user agent", map[string]any{"userAgent": "not found"})
		}

		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		s.logger.Info("request", map[string]any{
			"ip": ip,
			//"time":       time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			"method":     r.Method,
			"path":       r.URL.Path,
			"proto":      r.Proto,
			"status":     rw.status,
			"latency":    time.Since(start).Milliseconds(),
			"user_agent": userAgent,
		})
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *AppServer) helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (s *AppServer) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 Method Not Allowed"))
}
