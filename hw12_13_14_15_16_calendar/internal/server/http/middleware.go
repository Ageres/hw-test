package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
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

func (s *httpServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := handlehandleRequestID(r)
		logger := s.logger.With(map[string]any{
			"requestId":  utils.GetRequestID(ctx),
			"restMethod": r.Method,
		})
		ctx = logger.SetLoggerToCtx(ctx)

		ip := getIP(r, logger)
		userAgent := getUserAgent(r, logger)

		newR := r.WithContext(ctx)

		rw := &responseWriter{w, http.StatusOK}
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Header().Set(utils.RequestIDHeader, utils.GetRequestID(ctx))

		next.ServeHTTP(rw, newR)

		logger.Info("rest request", map[string]any{
			"ip":         ip,
			"path":       r.URL.Path,
			"protocol":   r.Proto,
			"status":     rw.status,
			"latency_ms": time.Since(start).Milliseconds(),
			"user_agent": userAgent,
		})
	})
}

func handlehandleRequestID(r *http.Request) context.Context {
	ctx := r.Context()
	requestId := r.Header.Get(utils.RequestIDHeader)
	if requestId == "" {
		ctx = utils.SetNewRequestIDToCtx(ctx)
	} else {
		ctx = utils.SetRequestIdToCtx(ctx, requestId)
	}
	return ctx
}

func getIP(r *http.Request, logger lg.Logger) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
		logger.WithError(err).Warn("get remote ip address", map[string]any{
			"ip": ip,
		})
		return ip
	}
	return ip
}

func getUserAgent(r *http.Request, logger lg.Logger) string {
	userAgent := r.UserAgent()
	if userAgent == "" {
		userAgent = "-"
		logger.Warn("get user agent", map[string]any{"userAgent": "not found"})
		return userAgent
	}
	return userAgent
}
