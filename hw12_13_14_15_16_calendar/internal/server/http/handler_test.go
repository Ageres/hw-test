package internalhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MockService struct{}

func (m *MockService) GetEventList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"test"}`))
}

func (m *MockService) AddEvent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"added"}`))
}

func (m *MockService) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}`))
}

func (m *MockService) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}`))
}

func (m *MockService) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"error":"method not allowed"}`))
}

func TestEventHandler(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)

	config := &model.HTTPConf{
		Server: &model.HTTPServerConf{
			ReadHeaderTimeout: 5,
			ReadTimeout:       10,
			WriteTimeout:      10,
			IdleTimeout:       30,
		},
	}

	service := &MockService{}
	server := NewHTTPServer(ctx, config, service).(*httpServer)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET method",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST method",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT method",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE method",
			method:         http.MethodDelete,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PATCH method not allowed",
			method:         http.MethodPatch,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/event", nil)
			w := httptest.NewRecorder()

			server.eventHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
