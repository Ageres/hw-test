package storage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEvent_GenerateEventId(t *testing.T) {
	event := &storage.Event{}
	require.Empty(t, event.ID)

	event.GenerateEventID()
	require.NotEmpty(t, event.ID)

	err := uuid.Validate(event.ID)
	require.NoError(t, err)
}

func TestValidateEventId(t *testing.T) {
	testCases := []struct {
		name     string
		id       string
		expected error
	}{
		{
			name:     "valid uuid",
			id:       uuid.New().String(),
			expected: nil,
		},
		{
			name:     "invalid uuid",
			id:       "invalid",
			expected: storage.NewSError("failed to validate event id", fmt.Errorf("invalid UUID length: 7")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := storage.ValidateEventID(tc.id)
			if tc.expected == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expected.Error())
			}
		})
	}
}

func TestValidateEvent(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(2 * time.Hour)
	pastTime := now.Add(-2 * time.Hour)

	testCases := []struct {
		name     string
		event    *storage.Event
		expected error
	}{
		{
			name: "valid event",
			event: &storage.Event{
				Title:     "Valid Event",
				StartTime: futureTime,
				Duration:  time.Hour,
				UserID:    "user1",
			},
			expected: nil,
		},
		{
			name:     "nil event",
			event:    nil,
			expected: storage.ErrEventIsNil,
		},
		{
			name: "empty title",
			event: &storage.Event{
				Title:     "",
				StartTime: futureTime,
				Duration:  time.Hour,
				UserID:    "user1",
			},
			expected: storage.NewSErrorWithMsgArr([]string{"title is empty"}),
		},
		{
			name: "expired event time",
			event: &storage.Event{
				Title:     "Past Event",
				StartTime: pastTime,
				Duration:  time.Hour,
				UserID:    "user1",
			},
			expected: storage.NewSErrorWithMsgArr([]string{"event time is expired"}),
		},
		{
			name: "invalid duration",
			event: &storage.Event{
				Title:     "Invalid Duration",
				StartTime: futureTime,
				Duration:  0,
				UserID:    "user1",
			},
			expected: storage.NewSErrorWithMsgArr([]string{"duration must be positive"}),
		},
		{
			name: "empty user id",
			event: &storage.Event{
				Title:     "No User",
				StartTime: futureTime,
				Duration:  time.Hour,
				UserID:    "",
			},
			expected: storage.NewSErrorWithMsgArr([]string{"user id is empty"}),
		},
		{
			name: "multiple errors",
			event: &storage.Event{
				Title:     "",
				StartTime: pastTime,
				Duration:  0,
				UserID:    "",
			},
			expected: storage.NewSErrorWithMsgArr([]string{
				"title is empty",
				"event time is expired",
				"duration must be positive",
				"user id is empty",
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := storage.ValidateEvent(tc.event)
			if tc.expected == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.expected.Error(), err.Error())
			}
		})
	}
}

func TestFullValidateEvent(t *testing.T) {
	now := time.Now()
	validID := uuid.New().String()

	testCases := []struct {
		name     string
		event    *storage.Event
		expected error
	}{
		{
			name: "fully valid event",
			event: &storage.Event{
				ID:        validID,
				Title:     "Valid Event",
				StartTime: now.Add(time.Hour),
				Duration:  time.Hour,
				UserID:    "user1",
			},
			expected: nil,
		},
		{
			name: "invalid uuid",
			event: &storage.Event{
				ID:        "invalid-uuid",
				Title:     "Test Event",
				StartTime: now.Add(time.Hour),
				Duration:  time.Hour,
				UserID:    "user1",
			},
			expected: storage.NewSErrorWithMsgArr([]string{
				fmt.Sprintf("failed to validate event id: %v", fmt.Errorf("invalid UUID length: 12")),
			}),
		},
		{
			name: "invalid uuid and other errors",
			event: &storage.Event{
				ID:        "invalid",
				Title:     "",
				StartTime: now.Add(-time.Hour),
				Duration:  0,
				UserID:    "",
			},
			expected: storage.NewSErrorWithMsgArr([]string{
				fmt.Sprintf("failed to validate event id: %v", fmt.Errorf("invalid UUID length: 7")),
				"title is empty",
				"event time is expired",
				"duration must be positive",
				"user id is empty",
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := storage.FullValidateEvent(tc.event)
			if tc.expected == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.expected.Error(), err.Error())
			}
		})
	}
}
