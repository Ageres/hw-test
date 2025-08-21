package internalhttp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

func (s *AppServer) helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (s *AppServer) methodNotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

type Period string

const (
	DAY   Period = "day"
	WEEK  Period = "week"
	MONTH Period = "month"
)

type ListRequest struct {
	Period   Period     `json:"period" binding:"required"`
	StartDay *time.Time `json:"startDay" binding:"required"`
}

func (s *AppServer) eventHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpMethod := r.Method
	switch httpMethod {
	case http.MethodGet:
		listRequest, err := unmarshalRequestBody[ListRequest](w, r)
		if err != nil {
			s.logger.WithError(err).Error("unmarshal get request body")
			return
		}
		switch listRequest.Period {
		case DAY:
			resp, err := s.app.ListDayEvents(ctx, *listRequest.StartDay)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			writeResponse(w, resp)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			return
		case WEEK:
			resp, err := s.app.ListWeekEvents(ctx, *listRequest.StartDay)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			writeResponse(w, resp)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			return
		case MONTH:
			resp, err := s.app.ListMonthEvents(ctx, *listRequest.StartDay)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			writeResponse(w, resp)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			return
		default:
			http.Error(w, "unknown period", http.StatusBadRequest)
			return
		}
	case http.MethodPost:
		req, err := unmarshalRequestBody[storage.Event](w, r)
		if err != nil {
			s.logger.WithError(err).Error("unmarshal post request body")
			return
		}
		resp, err := s.app.AddEvent(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeResponse(w, resp)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case http.MethodPut:
		req, err := unmarshalRequestBody[storage.Event](w, r)
		if err != nil {
			s.logger.WithError(err).Error("unmarshal get request body")
			return
		}
		resp, err := s.app.UpdateEvent(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeResponse(w, resp)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case http.MethodDelete:
		args := r.URL.Query()
		eventId := args.Get("eventId")
		resp, err := s.app.DeleteEvent(ctx, eventId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeResponse(w, resp)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
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
