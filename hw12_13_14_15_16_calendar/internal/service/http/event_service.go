package httpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type HttpService interface {
	GetEventList(w http.ResponseWriter, r *http.Request)
	AddEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	MethodNotAllowed(w http.ResponseWriter, r *http.Request)
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

func (h *httpService) AddEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := unmarshalRequestBody[AddEventRequest](ctx, w, r)
	if err != nil {
		return
	}
	resp, err := h.storage.Add(ctx, (*storage.Event)(req))
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("add event: %s", err.Error()),
			w,
			defineHttpStatusCode(err.Error()),
		)
		return
	}
	writeResponse(ctx, w, resp)
}

func (h *httpService) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := unmarshalRequestBody[UpdateEventRequest](ctx, w, r)
	if err != nil {
		return
	}
	err = h.storage.Update(ctx, (*storage.Event)(req))
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("update event: %s", err.Error()),
			w,
			defineHttpStatusCode(err.Error()),
		)
		return
	}
	resp := UpdateEventResponse{
		Status: UPDATE,
	}
	writeResponse(ctx, w, &resp)
}

func (h *httpService) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := unmarshalRequestBody[DeleteEventRequest](ctx, w, r)
	if err != nil {
		return
	}
	err = h.storage.Delete(ctx, req.Id)
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("delete event: %s", err.Error()),
			w,
			defineHttpStatusCode(err.Error()),
		)
		return
	}
	resp := DeleteEventResponse{
		Status: Delete,
	}
	writeResponse(ctx, w, &resp)
}

func (h *httpService) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeError(
		r.Context(),
		"Method Not Allowed",
		w,
		http.StatusMethodNotAllowed,
	)
}

//-----------------------------------------------------------------------------------
// вспомогательные функции

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
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("write responce error: %s", err.Error()),
			w,
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeError(ctx context.Context, errMsg string, w http.ResponseWriter, httpSatus int) {
	httpError := NewHttpError(ctx, errMsg)
	err := json.NewEncoder(w).Encode(httpError)
	if err != nil {
		lg.GetLogger(ctx).WithError(err).Error("encode http error", map[string]any{"errMsg": errMsg})
		errMsg := fmt.Sprintf("encode http error: errMsg '%s', error '%s'", errMsg, err.Error())
		errResp := fmt.Sprintf("{\"error\": \"%s\"}", errMsg)
		http.Error(w, errResp, httpSatus)
		return
	}
	lg.GetLogger(ctx).WithError(httpError).Error("write error")
}

func defineHttpStatusCode(errMsg string) int {
	if strings.Contains(errMsg, "user is not the owner of the event, conflict with") || strings.Contains(errMsg, "time is already taken by another event") {
		return http.StatusConflict
	}
	if strings.Contains(errMsg, "event not found") {
		return http.StatusNotFound
	}
	if strings.Contains(errMsg, "failed to validate event id") ||
		strings.Contains(errMsg, "title is empty") ||
		strings.Contains(errMsg, "event time is expired") ||
		strings.Contains(errMsg, "duration must be positive") ||
		strings.Contains(errMsg, "user id is empty") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
