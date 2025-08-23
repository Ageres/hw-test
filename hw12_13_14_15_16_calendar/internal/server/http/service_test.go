package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	bserv "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/http/baseserver"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Add(ctx context.Context, event *storage.Event) (*storage.Event, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.Event), args.Error(1)
}

func (m *MockStorage) Update(ctx context.Context, event *storage.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockStorage) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorage) ListDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, startDay)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) ListWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, startDay)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) ListMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, startDay)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func TestHttpService_GetEventList_Ok(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	timeLocation, _ := time.LoadLocation("UTC")
	start := time.Date(2030, 12, 31, 0, 0, 0, 0, timeLocation)

	testEvents := []storage.Event{
		{ID: "event-1", Title: "Test Event 1"},
		{ID: "event-2", Title: "Test Event 2"},
	}

	tests := []struct {
		name           string
		method         string
		requestBody    any
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful day events",
			method: http.MethodGet,
			requestBody: bserv.GetEventListRequest{
				Period:   bserv.DAY,
				StartDay: &start,
			},
			mockSetup: func() {
				mockStorage.On("ListDay", mock.Anything, start).Return(testEvents, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"Day event list successfully retrieved"`,
		},
		{
			name:   "successful week events",
			method: http.MethodGet,
			requestBody: bserv.GetEventListRequest{
				Period:   bserv.WEEK,
				StartDay: &start,
			},
			mockSetup: func() {
				mockStorage.On("ListWeek", mock.Anything, start).Return(testEvents, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"Week event list successfully retrieved"`,
		},
		{
			name:   "successful month events",
			method: http.MethodGet,
			requestBody: bserv.GetEventListRequest{
				Period:   bserv.MONTH,
				StartDay: &start,
			},
			mockSetup: func() {
				mockStorage.On("ListMonth", mock.Anything, start).Return(testEvents, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"Month event list successfully retrieved"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(tt.method, "/v1/event", bytes.NewReader(body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			service.GetEventList(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestHttpService_MethodNotAllowed_Error(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	req := httptest.NewRequest(http.MethodPatch, "/v1/event", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	service.MethodNotAllowed(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Contains(t, w.Body.String(), `"Method Not Allowed"`)
}

func TestDefineHttpStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected int
	}{
		{
			name:     "conflict - user not owner",
			errMsg:   "user is not the owner of the event, conflict with",
			expected: http.StatusConflict,
		},
		{
			name:     "conflict - time taken",
			errMsg:   "time is already taken by another event",
			expected: http.StatusConflict,
		},
		{
			name:     "not found",
			errMsg:   "event not found",
			expected: http.StatusNotFound,
		},
		{
			name:     "bad request - validation",
			errMsg:   "failed to validate event id",
			expected: http.StatusBadRequest,
		},
		{
			name:     "bad request - title empty",
			errMsg:   "title is empty",
			expected: http.StatusBadRequest,
		},
		{
			name:     "internal server error",
			errMsg:   "unknown error",
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defineHttpStatusCode(tt.errMsg)
			assert.Equal(t, tt.expected, result)
		})
	}
}
