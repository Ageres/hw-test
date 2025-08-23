package internalhttp

import (
	"context"
	"net/http"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	bserv "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http/baseserver"
)

type httpServer struct {
	server  *http.Server
	logger  lg.Logger
	address string
	service HTTPService
}

func NewHTTPServer(ctx context.Context, httpConf *model.HTTPConf, service HTTPService) bserv.HTTPServer {
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
		service: service,
	}

	s.server.Handler = s.createRouter()

	s.logger.Info("server configured")
	return s
}

func (s *httpServer) createRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	mux.HandleFunc("/v1/event", s.eventHandler)
	mux.HandleFunc("/", s.methodNotAllowedHandler)
	return s.loggingMiddleware(mux)
}

func (s *httpServer) Start(_ context.Context) error {
	s.logger.Info("Starting HTTP server", map[string]any{
		"address": s.address,
	})
	return s.server.ListenAndServe()
}

func (s *httpServer) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
