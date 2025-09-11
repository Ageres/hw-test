package internalhttp

import (
	"context"
	"net/http"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type AppServer struct {
	server  *http.Server
	logger  lg.Logger
	app     Application
	address string
}

type Application interface { // TODO
}

func NewServer(ctx context.Context, httpConf *model.HTTPConf, app Application) Server {
	address := httpConf.Server.GetAddress()

	s := &AppServer{
		server: &http.Server{
			Addr:              address,
			ReadHeaderTimeout: time.Duration(httpConf.Server.ReadHeaderTimeout) * time.Second,
			ReadTimeout:       time.Duration(httpConf.Server.ReadTimeout) * time.Second,
			WriteTimeout:      time.Duration(httpConf.Server.WriteTimeout) * time.Second,
			IdleTimeout:       time.Duration(httpConf.Server.IdleTimeout) * time.Second,
		},
		logger:  lg.GetLogger(ctx),
		app:     app,
		address: address,
	}

	s.server.Handler = s.createRouter()

	s.logger.Info("server configured")
	return s
}

func (s *AppServer) createRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	mux.HandleFunc("/", s.methodNotAllowedHandler)
	return s.loggingMiddleware(mux)
}

func (s *AppServer) Start(_ context.Context) error {
	s.logger.Info("Starting HTTP server", map[string]any{
		"address": s.address,
	})
	return s.server.ListenAndServe()
}

func (s *AppServer) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
