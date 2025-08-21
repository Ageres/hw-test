package httpservice

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type HttpService interface {
	GetEventList(w http.ResponseWriter, r *http.Request)
	AddEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
}

type httpService struct {
	storage storage.Storage
}

func NewHttpService(ctx context.Context, storage storage.Storage) HttpService {
	return &httpService{
		storage: storage,
	}
}

// GetEventList implements HttpServece.
func (h *httpService) GetEventList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := lg.GetLogger(ctx)
	listRequest, err := unmarshalRequestBody[GetEventListRequest](w, r)
	if err != nil {
		logger.WithError(err).Error("unmarshal get request body")
		return
	}
	switch listRequest.Period {
	case DAY:
		h.getEventList(ctx, w, listRequest.StartDay, LISTDAY, h.storage.ListDay)
		return
	case WEEK:
		h.getEventList(ctx, w, listRequest.StartDay, LISTWEEK, h.storage.ListWeek)
		return
	case MONTH:
		h.getEventList(ctx, w, listRequest.StartDay, LISTMONTH, h.storage.ListMonth)
		return
	default:
		http.Error(w, "unknown period", http.StatusBadRequest)
		return
	}
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

func unmarshalRequestBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	reqRef := new(T)
	err = json.Unmarshal(buf, reqRef)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return reqRef, nil
}

func (h *httpService) getEventList(
	ctx context.Context,
	w http.ResponseWriter,
	startDay *time.Time,
	status GetEventListStatus,
	list func(ctx context.Context, startDay time.Time) ([]storage.Event, error),
) {
	if startDay == nil {
		http.Error(w, "startDay is nil", http.StatusBadRequest)
		return
	}
	events, err := list(ctx, *startDay)
	if err != nil {
		lg.GetLogger(ctx).WithError(err).Error("get event list")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	resp := GetListResponse{
		Status: status,
		Events: events,
	}
	writeResponse(w, &resp)
}

func writeResponse[T any](w http.ResponseWriter, resp *T) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		slog.Error("responce marshal error", "err", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		slog.Error("responce marshal error", "err", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
