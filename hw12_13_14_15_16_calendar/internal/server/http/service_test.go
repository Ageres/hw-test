package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (m *MockStorage) ListReminderEvents(ctx context.Context, startTime, endTime time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) ResetEventReminder(ctx context.Context, eventIDs []string) error {
	args := m.Called(ctx, eventIDs)
	return args.Error(0)
}

func (m *MockStorage) Close() error {
	return nil
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

func TestHttpService_GetEventList_Error(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	timeLocation, _ := time.LoadLocation("UTC")
	start := time.Date(2030, 12, 31, 0, 0, 0, 0, timeLocation)

	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "invalid period",
			method: http.MethodGet,
			requestBody: map[string]interface{}{
				"period":   "invalid",
				"startDay": start.Format(time.RFC3339),
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"unknown period: invalid"`,
		},
		{
			name:   "missing startDay",
			method: http.MethodGet,
			requestBody: map[string]interface{}{
				"period": "day",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"startDay is nil"`,
		},
		{
			name:   "storage error",
			method: http.MethodGet,
			requestBody: bserv.GetEventListRequest{
				Period:   bserv.DAY,
				StartDay: &start,
			},
			mockSetup: func() {
				mockStorage.On("ListDay", mock.Anything, start).Return([]storage.Event{}, fmt.Errorf("storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"get event list: storage error"`,
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

func TestHttpService_AddEvent_Ok(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	timeLocation, _ := time.LoadLocation("UTC")
	start := time.Date(2030, 12, 31, 0, 0, 0, 0, timeLocation)

	event := &storage.Event{
		Title:     "Test Event",
		StartTime: start,
		Duration:  time.Hour,
		UserID:    "user-1",
	}

	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful add",
			requestBody: bserv.AddEventRequest(*event),
			mockSetup: func() {
				returnedEvent := *event
				returnedEvent.ID = "new-id"
				mockStorage.On("Add", mock.Anything, mock.AnythingOfType("*storage.Event")).
					Return(&returnedEvent, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"Test Event"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/event", bytes.NewReader(body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			service.AddEvent(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestHttpService_AddEvent_Error(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	timeLocation, _ := time.LoadLocation("UTC")
	start := time.Date(2030, 12, 31, 0, 0, 0, 0, timeLocation)

	event := &storage.Event{
		Title:     "Test Event",
		StartTime: start,
		Duration:  time.Hour,
		UserID:    "user-1",
	}

	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "storage error",
			requestBody: bserv.AddEventRequest(*event),
			mockSetup: func() {
				mockStorage.On("Add", mock.Anything, mock.AnythingOfType("*storage.Event")).
					Return(nil, fmt.Errorf("storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "add event: storage error",
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unmarshal request body: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/event", bytes.NewReader(body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			service.AddEvent(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestHttpService_UpdateEvent_Ok(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	event := &storage.Event{
		ID:    "test-id",
		Title: "Updated Event",
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful update",
			requestBody: bserv.UpdateEventRequest(*event),
			mockSetup: func() {
				mockStorage.On("Update", mock.Anything, mock.AnythingOfType("*storage.Event")).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"Event updated successfully"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/event", bytes.NewReader(body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			service.UpdateEvent(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestHttpService_UpdateEvent_Error(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	event := &storage.Event{
		ID:    "test-id",
		Title: "Updated Event",
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "storage error",
			requestBody: bserv.UpdateEventRequest(*event),
			mockSetup: func() {
				mockStorage.On("Update", mock.Anything, mock.AnythingOfType("*storage.Event")).
					Return(fmt.Errorf("storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "update event: storage error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/event", bytes.NewReader(body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			service.UpdateEvent(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestHttpService_DeleteEvent_Ok(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful delete",
			requestBody: bserv.DeleteEventRequest{
				ID: "test-id",
			},
			mockSetup: func() {
				mockStorage.On("Delete", mock.Anything, "test-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"Event deleted successfully"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/event", bytes.NewReader(body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			service.DeleteEvent(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestHttpService_DeleteEvent_Error(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockStorage := new(MockStorage)

	service := &httpService{
		storage: mockStorage,
	}

	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "storage error",
			requestBody: bserv.DeleteEventRequest{
				ID: "test-id",
			},
			mockSetup: func() {
				mockStorage.On("Delete", mock.Anything, "test-id").Return(fmt.Errorf("storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "delete event: storage error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/event", bytes.NewReader(body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			service.DeleteEvent(w, req)

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
