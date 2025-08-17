package internalhttp

import (
	"context"
	"net/http"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type AppServer struct {
	httpServer *http.Server
	app        Application
	logger     lg.Logger
}

type Application interface { // TODO
}

type AppHandler struct {
	Method  string                                       `json:"method"`
	Path    string                                       `json:"path"`
	Handler func(w http.ResponseWriter, r *http.Request) `json:"-"`
}

func NewServer(ctx context.Context, httpConf *model.HttpConf, app Application) Server {
	address := httpConf.Server.GetAddress()
	loggger := lg.GetLogger(ctx)
	appServer := AppServer{
		httpServer: &http.Server{
			Addr: address,
		},
		app:    app,
		logger: loggger,
	}
	server := appServer.configureMux()
	logger.GetLogger(ctx).Info("server configured", map[string]any{"address": address})
	return server
}

func (s *AppServer) Start(ctx context.Context) error {
	return s.httpServer.ListenAndServe()
}

func (s *AppServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *AppServer) configureMux() Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	mux.HandleFunc("/", s.methodNotAllowed)
	s.httpServer.Handler = s.loggingMiddleware(mux)
	return s
}
