package storage_test

import (
	"testing"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestEvent_ToNotification_NilEvent(t *testing.T) {
	var event *storage.Event
	notification := event.ToNotification()
	require.Nil(t, notification)
}

func TestEvent_ToNotification(t *testing.T) {
	now := time.Now()
	event := &storage.Event{
		ID:        "test-id",
		Title:     "Test Event",
		StartTime: now,
		UserID:    "user1",
	}

	notification := event.ToNotification()
	require.NotNil(t, notification)
	require.Equal(t, event.ID, notification.ID)
	require.Equal(t, event.Title, notification.Title)
	require.Equal(t, event.UserID, notification.UserID)
	require.Equal(t, now.Year(), notification.Date.Year())
	require.Equal(t, now.Month(), notification.Date.Month())
	require.Equal(t, now.Day(), notification.Date.Day())
	require.Zero(t, notification.Date.Hour())
	require.Zero(t, notification.Date.Minute())
	require.Zero(t, notification.Date.Second())
}
