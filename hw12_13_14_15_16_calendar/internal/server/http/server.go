package internalhttp

import (
	"context"
	"net/http"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type MyServer struct {
	httpServer *http.Server
	app        Application
}

type Logger interface { // TODO
	//code
}

type Application interface { // TODO
}

type Myhandler struct {
	Method  string                                       `json:"method"`
	Path    string                                       `json:"path"`
	Handler func(w http.ResponseWriter, r *http.Request) `json:"-"`
}

var myHandlers = []Myhandler{
	{"get", "/hello", helloHandler},
	{"any other", "/", notFoundHandler},
}

func NewServer(ctx context.Context, httpConf *model.HttpConf, app Application) Server {
	address := httpConf.Server.GetAddress()
	mux := configureMux()

	server := MyServer{
		httpServer: &http.Server{
			Addr:    address,
			Handler: loggingMiddleware(mux),
		},
		app: app,
	}
	logger.GetLogger(ctx).Info("server configured", map[string]any{
		"Addr":     address,
		"handlers": myHandlers,
	})
	return &server
}

func (s *MyServer) Start(ctx context.Context) error {
	return s.httpServer.ListenAndServe()
}

func (s *MyServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func configureMux() *http.ServeMux {
	mux := http.NewServeMux()
	for _, handler := range myHandlers {
		mux.HandleFunc(handler.Path, handler.Handler)
	}
	return mux
}
