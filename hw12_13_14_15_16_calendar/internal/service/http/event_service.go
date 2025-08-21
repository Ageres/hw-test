package httpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (h *httpService) GetEventList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	listRequest, err := unmarshalRequestBody[GetEventListRequest](ctx, w, r)
	if err != nil {
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
		writeError(
			ctx,
			fmt.Sprintf("unknown period: %s", listRequest.Period),
			w,
			http.StatusBadRequest,
		)
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

func unmarshalRequestBody[T any](ctx context.Context, w http.ResponseWriter, r *http.Request) (*T, error) {
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		writeError(
			ctx,
			fmt.Sprintf("read request body: %s", err.Error()),
			w,
			http.StatusBadRequest,
		)
		return nil, err
	}
	reqRef := new(T)
	err = json.Unmarshal(buf, reqRef)
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("unmarshal request body: %s", err.Error()),
			w,
			http.StatusBadRequest,
		)
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
		writeError(ctx, "startDay is nil", w, http.StatusBadRequest)
		return
	}
	events, err := list(ctx, *startDay)
	if err != nil {
		lg.GetLogger(ctx).WithError(err).Error("get event list")
		writeError(
			ctx,
			fmt.Sprintf("get event list: %s", err.Error()),
			w,
			http.StatusInternalServerError,
		)
		return
	}
	resp := GetListResponse{
		Status: status,
		Events: events,
	}
	writeResponse(ctx, w, &resp)
}

func writeResponse[T any](ctx context.Context, w http.ResponseWriter, resp *T) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("responce marshal error: %s", err.Error()),
			w,
			http.StatusInternalServerError,
		)
		return
	}
	_, err = w.Write(resBuf)
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("write responce error: %s", err.Error()),
			w,
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func writeError(ctx context.Context, errMsg string, w http.ResponseWriter, httpSatus int) {
	httpError := NewHttpError(errMsg)
	data, err := json.Marshal(httpError)
	if err != nil {
		lg.GetLogger(ctx).WithError(err).Error("marshal http error", map[string]any{"errMsg": errMsg})
		return
	}
	lg.GetLogger(ctx).WithError(httpError).Error("write error")
	http.Error(w, string(data), httpSatus)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
