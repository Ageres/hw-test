package internalhttp

import (
	"context"
	"net/http"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	bserv "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http/baseserver"
)

func NewHealthCheckHTTPServer(ctx context.Context, httpConf *model.HealthHTTPConf) bserv.HTTPServer {
	address := httpConf.Server.GetAddress()

	s := &httpServer{
		server: &http.Server{
			Addr:              address,
			ReadHeaderTimeout: time.Duration(httpConf.Server.ReadHeaderTimeout) * time.Second,
			ReadTimeout:       time.Duration(httpConf.Server.ReadTimeout) * time.Second,
			WriteTimeout:      time.Duration(httpConf.Server.WriteTimeout) * time.Second,
			IdleTimeout:       time.Duration(httpConf.Server.IdleTimeout) * time.Second,
		},
		logger:  lg.GetLogger(ctx),
		address: address,
	}

	s.server.Handler = s.createHealtCheckhRouter()

	s.logger.Info("health check server configured")
	return s
}

func (s *httpServer) createHealtCheckhRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthCheckHandler)
	mux.HandleFunc("/", s.methodNotAllowedHandler)
	return s.loggingMiddleware(mux)
}
