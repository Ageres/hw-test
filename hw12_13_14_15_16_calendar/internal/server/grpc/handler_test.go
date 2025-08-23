package internalgrpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
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

func TestGrpcServer_GetEvent_Ok(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockLogger := lg.GetLogger(ctx)
	mockStorage := new(MockStorage)

	server := &GrpcServer{
		storage: mockStorage,
		logger:  mockLogger,
	}

	startTime := timestamppb.New(time.Now())
	validRequest := &pb.GetEventListRequest{
		Period:    pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_DAY,
		StartTime: startTime,
	}

	tests := []struct {
		name          string
		request       *pb.GetEventListRequest
		mockSetup     func()
		expectedError bool
		errorCode     codes.Code
		errorMessage  string
	}{
		{
			name:    "successful day events",
			request: validRequest,
			mockSetup: func() {
				mockStorage.On("ListDay", mock.Anything, startTime.AsTime()).
					Return([]storage.Event{{ID: "test-id"}}, nil)
			},
			expectedError: false,
		},
		{
			name:    "successful week events",
			request: &pb.GetEventListRequest{Period: pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_WEEK, StartTime: startTime},
			mockSetup: func() {
				mockStorage.On("ListWeek", mock.Anything, startTime.AsTime()).
					Return([]storage.Event{{ID: "test-id"}}, nil)
			},
			expectedError: false,
		},
		{
			name:    "successful month events",
			request: &pb.GetEventListRequest{Period: pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_MONTH, StartTime: startTime},
			mockSetup: func() {
				mockStorage.On("ListMonth", mock.Anything, startTime.AsTime()).
					Return([]storage.Event{{ID: "test-id"}}, nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := server.GetEvent(ctx, tt.request)

			if tt.expectedError {
				require.Error(t, err)
				if tt.errorCode != 0 {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.errorCode, st.Code())
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Len(t, resp.Events, 1)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestGrpcServer_GetEvent_Error(t *testing.T) {
	ctx := utils.SetNewRequestIDToCtx(context.Background())
	ctx = lg.SetDefaultLogger(ctx)
	mockLogger := lg.GetLogger(ctx)
	mockStorage := new(MockStorage)

	server := &GrpcServer{
		storage: mockStorage,
		logger:  mockLogger,
	}

	startTime := timestamppb.New(time.Now())
	validRequest := &pb.GetEventListRequest{
		Period:    pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_DAY,
		StartTime: startTime,
	}

	tests := []struct {
		name          string
		request       *pb.GetEventListRequest
		mockSetup     func()
		expectedError bool
		errorCode     codes.Code
		errorMessage  string
	}{
		{
			name:    "invalid period",
			request: &pb.GetEventListRequest{Period: pb.GetEventListPeriod_GET_EVENT_LIST_PERIOD_UNSPECIFIED, StartTime: startTime},
			mockSetup: func() {
			},
			expectedError: true,
			errorCode:     codes.InvalidArgument,
			errorMessage:  "invalid period",
		},
		{
			name:    "storage error",
			request: validRequest,
			mockSetup: func() {
				mockStorage.On("ListDay", mock.Anything, startTime.AsTime()).
					Return([]storage.Event{}, errors.New("storage error"))
			},
			expectedError: true,
			errorCode:     codes.Internal,
			errorMessage:  "storage error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := server.GetEvent(ctx, tt.request)

			if tt.expectedError {
				require.Error(t, err)
				if tt.errorCode != 0 {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.errorCode, st.Code())
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Len(t, resp.Events, 1)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
