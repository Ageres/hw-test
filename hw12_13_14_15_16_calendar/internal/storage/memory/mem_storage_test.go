package memorystorage

import (
	"testing"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage_Add(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("add events", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[2]))
		require.Len(t, dto.storage.events, 3)
	})

	t.Run("add nil event error when adding", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))

		err := dto.storage.Add(nil)
		require.ErrorIs(t, err, storage.ErrEventIsNil)
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("validation event error when adding", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))

		err := dto.storage.Add(&events[4])
		require.Error(t, err)
		require.Equal(t, err.Error(), "title is empty, event time is expired, user id is empty")
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("event duplication error when adding", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))

		err := dto.storage.Add(&events[0])
		require.ErrorIs(t, err, storage.ErrEventAllreadyExists)
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("date busy error when adding", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))

		err := dto.storage.Add(&events[3])
		require.ErrorIs(t, err, storage.ErrDateBusy)
		require.Len(t, dto.storage.events, 1)
	})

}

type TestMemoryStorageDto struct {
	storage *MemoryStorage
	events  []storage.Event
}

func newTestMemoryStorageDto() *TestMemoryStorageDto {
	return &TestMemoryStorageDto{}
}

func (dto *TestMemoryStorageDto) buildNewStorage() *TestMemoryStorageDto {
	dto.storage = NewMemoryStorage().(*MemoryStorage)
	return dto
}

func (dto *TestMemoryStorageDto) buildNewEvents() *TestMemoryStorageDto {
	now := time.Now()
	userID := uuid.New().String()
	event0 := storage.Event{
		ID:          uuid.New().String(),
		Title:       "Event 0",
		StartTime:   now.Add(1 * time.Hour),
		Duration:    30 * time.Minute,
		Description: "Test event 0",
		UserID:      userID,
	}
	event1 := storage.Event{
		Title:     "Event 1",
		StartTime: now.Add(2 * time.Hour),
		Duration:  1 * time.Hour,
		UserID:    userID,
	}
	event2 := storage.Event{
		Title:     "Event 2 (other user)",
		StartTime: now.Add(1 * time.Hour),
		Duration:  30 * time.Minute,
		UserID:    uuid.New().String(),
	}
	event3 := storage.Event{
		Title:       "Event 3",
		StartTime:   now.Add(50 * time.Minute),
		Duration:    30 * time.Minute,
		Description: "Test event 3",
		UserID:      userID,
	}
	event4 := storage.Event{
		StartTime:   now.Add(-50 * time.Minute),
		Duration:    30 * time.Minute,
		Description: "Test event 4",
	}
	dto.events = append(dto.events, event0, event1, event2, event3, event4)
	return dto
}
