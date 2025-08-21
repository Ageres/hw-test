package internalhttp

import (
	"net/http"
)

func (s *AppServer) helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.service.MethodNotAllowed(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (s *AppServer) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	s.service.MethodNotAllowed(w, r)
}

func (s *AppServer) eventHandler(w http.ResponseWriter, r *http.Request) {
	httpMethod := r.Method
	switch httpMethod {
	case http.MethodGet:
		s.service.GetEventList(w, r)
		return
	case http.MethodPost:
		s.service.AddEvent(w, r)
		return
	case http.MethodPut:
		s.service.UpdateEvent(w, r)
		return
	case http.MethodDelete:
		s.service.DeleteEvent(w, r)
		return
	default:
		s.service.MethodNotAllowed(w, r)
		return
	}
}
