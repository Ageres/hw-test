package internalhttp

import (
	"context"
	"net/http"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Server struct {
	httpServer *http.Server
}

type Logger interface { // TODO
	//code
}

type Application interface { // TODO
}

func NewServer(ctx context.Context, httpConf *model.HttpConf, app Application) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:    httpConf.Server.GetAddress(),
			Handler: loggingMiddleware(mux),
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	//<-ctx.Done()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return s.httpServer.Shutdown(ctx)
}

// TODO
