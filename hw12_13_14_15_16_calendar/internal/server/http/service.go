package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	bserv "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http/baseserver"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
)

type HTTPService interface {
	GetEventList(w http.ResponseWriter, r *http.Request)
	AddEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	MethodNotAllowed(w http.ResponseWriter, r *http.Request)
}

type httpService struct {
	storage storage.Storage
}

func NewHTTPService(storage storage.Storage) HTTPService {
	return &httpService{
		storage: storage,
	}
}

func (h *httpService) GetEventList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	listRequest, err := unmarshalRequestBody[bserv.GetEventListRequest](ctx, w, r)
	if err != nil {
		return
	}
	switch listRequest.Period {
	case bserv.DAY:
		h.getEventList(ctx, w, listRequest.StartDay, bserv.LISTDAY, h.storage.ListDay)
		return
	case bserv.WEEK:
		h.getEventList(ctx, w, listRequest.StartDay, bserv.LISTWEEK, h.storage.ListWeek)
		return
	case bserv.MONTH:
		h.getEventList(ctx, w, listRequest.StartDay, bserv.LISTMONTH, h.storage.ListMonth)
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
	req, err := unmarshalRequestBody[bserv.AddEventRequest](ctx, w, r)
	if err != nil {
		return
	}
	resp, err := h.storage.Add(ctx, (*storage.Event)(req))
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("add event: %s", err.Error()),
			w,
			utils.DefineStatusCode(err.Error()),
		)
		return
	}
	writeResponse(ctx, w, resp)
}

func (h *httpService) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := unmarshalRequestBody[bserv.UpdateEventRequest](ctx, w, r)
	if err != nil {
		return
	}
	err = h.storage.Update(ctx, (*storage.Event)(req))
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("update event: %s", err.Error()),
			w,
			utils.DefineStatusCode(err.Error()),
		)
		return
	}
	resp := bserv.UpdateEventResponse{
		Status: bserv.UPDATE,
	}
	writeResponse(ctx, w, &resp)
}

func (h *httpService) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := unmarshalRequestBody[bserv.DeleteEventRequest](ctx, w, r)
	if err != nil {
		return
	}
	err = h.storage.Delete(ctx, req.ID)
	if err != nil {
		writeError(
			ctx,
			fmt.Sprintf("delete event: %s", err.Error()),
			w,
			utils.DefineStatusCode(err.Error()),
		)
		return
	}
	resp := bserv.DeleteEventResponse{
		Status: bserv.DELETE,
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
	if err != nil && err != io.EOF { //nolint:errorlint
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
	status bserv.GetEventListStatus,
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
	resp := bserv.GetEventListResponse{
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
			fmt.Sprintf("write response error: %s", err.Error()),
			w,
			http.StatusInternalServerError,
		)
		return
	}
}

func writeError(ctx context.Context, errMsg string, w http.ResponseWriter, httpSatus int) {
	w.WriteHeader(httpSatus)
	httpError := model.NewCalendarServiceError(httpSatus, errMsg, utils.GetRequestID(ctx), nil)
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
