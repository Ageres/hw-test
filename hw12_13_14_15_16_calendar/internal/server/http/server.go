package internalhttp

import (
	"context"
	"net/http"

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

func NewServer(ctx context.Context, httpConf *model.HttpConf, app Application) Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", notFoundHandler)
	return &MyServer{
		httpServer: &http.Server{
			Addr:    httpConf.Server.GetAddress(),
			Handler: loggingMiddleware(mux),
		},
		app: app,
	}
}

func (s *MyServer) Start(ctx context.Context) error {
	return s.httpServer.ListenAndServe()
}

func (s *MyServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
