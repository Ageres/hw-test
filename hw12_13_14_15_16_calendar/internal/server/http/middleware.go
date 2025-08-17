package internalhttp

import (
	"net/http"
)

type MyMiddleware struct {
}

func (s *AppServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
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
