package storage_test

import (
	"testing"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestEvent_ToNotification_NilEvent(t *testing.T) {
	var event *storage.Event
	notification := event.ToNotification()
	require.Nil(t, notification)
}
