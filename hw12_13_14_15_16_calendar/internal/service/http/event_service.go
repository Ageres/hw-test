package httpservice

import (
	"context"
	"net/http"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type HttpServece interface {
	GetEventList(w http.ResponseWriter, r *http.Request)
	AddEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
}

type httpService struct {
	storage storage.Storage
}

func NewHttpService(ctx context.Context, storage storage.Storage) HttpServece {
	return &httpService{
		storage: storage,
	}
}

// GetEventList implements HttpServece.
func (h *httpService) GetEventList(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// AddEvent implements HttpServece.
func (h *httpService) AddEvent(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// UpdateEvent implements HttpServece.
func (h *httpService) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// DeleteEvent implements HttpServece.
func (h *httpService) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
