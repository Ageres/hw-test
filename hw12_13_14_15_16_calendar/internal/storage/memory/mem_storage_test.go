package memorystorage

import (
	"testing"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorageAdd(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("add event", func(t *testing.T) {
		dto.buildNewStorage()
		require.Len(t, dto.storage.events, 0)
		require.NoError(t, dto.storage.Add(&events[0]))
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[2]))
		require.Len(t, dto.storage.events, 3)
	})

	t.Run("nil event error when adding", func(t *testing.T) {
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
		require.Equal(t, err.Error(), "title is empty; event time is expired; user id is empty")
		require.Len(t, dto.storage.events, 1)

		err = dto.storage.Add(&events[5])
		require.Error(t, err)
		require.Equal(t, err.Error(), "validate event id: invalid UUID length: 22; title is empty; event time is expired; user id is empty")
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

func TestStorageUpdate(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("update event", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))
		require.Len(t, dto.storage.events, 1)
		require.NoError(t, dto.storage.Update(events[0].ID, &events[1]))
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("nil event error when updating", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))

		err := dto.storage.Update(events[0].ID, nil)
		require.ErrorIs(t, err, storage.ErrEventIsNil)
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("validation event error when updating", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))

		err := dto.storage.Update("", &events[4])
		require.Error(t, err)
		require.Equal(t, err.Error(), "validate event id: invalid UUID length: 0; title is empty; event time is expired; user id is empty")
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("event not found error when updating", func(t *testing.T) {
		dto.buildNewStorage()

		err := dto.storage.Update(events[0].ID, &events[0])
		require.ErrorIs(t, err, storage.ErrEventNotFound)
		require.Len(t, dto.storage.events, 0)
	})

	t.Run("user conflict error when updating", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))
		require.Len(t, dto.storage.events, 1)

		err := dto.storage.Update(events[0].ID, &events[2])
		require.ErrorIs(t, err, storage.ErrUserConflict)
		require.Len(t, dto.storage.events, 1)
	})
}

func TestStorageDelete(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("delete event", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[2]))
		require.Len(t, dto.storage.events, 3)
		require.NoError(t, dto.storage.Delete(events[0].ID))
		require.Len(t, dto.storage.events, 2)
	})

	t.Run("validation id error when deleting", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))
		require.Len(t, dto.storage.events, 1)

		err := dto.storage.Delete("")
		require.Error(t, err)
		require.Equal(t, err.Error(), "validate event id: invalid UUID length: 0")
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("event not found error when deleting", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[2]))
		require.Len(t, dto.storage.events, 2)
		err := dto.storage.Delete(events[0].ID)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
		require.Len(t, dto.storage.events, 2)
	})
}

func TestStorageListEvents(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("list day events", func(t *testing.T) {
		dto.buildNewStorage()
		require.Len(t, dto.storage.events, 0)
		require.NoError(t, dto.storage.Add(&events[0]))
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[2]))
		require.Len(t, dto.storage.events, 3)

		dayEvents, err := dto.storage.ListDay(dto.now)
		require.NoError(t, err)
		require.Len(t, dayEvents, 3)

		require.NoError(t, dto.storage.Delete(events[0].ID))
		require.NoError(t, dto.storage.Delete(events[1].ID))
		require.Len(t, dto.storage.events, 1)

		dayEvents, err = dto.storage.ListDay(dto.now)
		require.NoError(t, err)
		require.Equal(t, dayEvents[0].ID, events[2].ID)
	})

	t.Run("list week events", func(t *testing.T) {
		dto.buildNewStorage()
		require.Len(t, dto.storage.events, 0)
		require.NoError(t, dto.storage.Add(&events[0]))
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[2]))
		require.Len(t, dto.storage.events, 3)

		dayEvents, err := dto.storage.ListDay(dto.now)
		require.NoError(t, err)
		require.Len(t, dayEvents, 3)

		require.NoError(t, dto.storage.Delete(events[0].ID))
		require.NoError(t, dto.storage.Delete(events[1].ID))
		require.Len(t, dto.storage.events, 1)

		dayEvents, err = dto.storage.ListDay(dto.now)
		require.NoError(t, err)
		require.Equal(t, dayEvents[0].ID, events[2].ID)
	})

}

// -------------------------------------------------------------------------------------
// Вспомогательные функции
type TestMemoryStorageDto struct {
	storage *MemoryStorage
	events  []storage.Event
	now     time.Time
}

func newTestMemoryStorageDto() *TestMemoryStorageDto {
	return &TestMemoryStorageDto{}
}

func (dto *TestMemoryStorageDto) buildNewStorage() *TestMemoryStorageDto {
	dto.storage = NewMemoryStorage().(*MemoryStorage)
	return dto
}

func (dto *TestMemoryStorageDto) buildNewEvents() *TestMemoryStorageDto {
	dto.now = time.Now()
	correctEventId := "aaeef68f-267d-459d-bda6-c900e27f4afe"
	wrongEventId := "459d-bda6-c900e27f4afe"
	userIDOne := "d6e2955f-7a5b-47f2-8f03-999ad489f51a"
	userIDTwo := "6cf51d87-ab61-437e-9c8c-193984d07bf6"
	event0 := storage.Event{
		ID:          correctEventId,
		Title:       "Event 0",
		StartTime:   dto.now.Add(1 * time.Hour),
		Duration:    30 * time.Minute,
		Description: "Test event 0",
		UserID:      userIDOne,
	}
	event1 := storage.Event{
		Title:     "Event 1",
		StartTime: dto.now.Add(2 * time.Hour),
		Duration:  1 * time.Hour,
		UserID:    userIDOne,
	}
	event2 := storage.Event{
		Title:     "Event 2 (other user)",
		StartTime: dto.now.Add(1 * time.Hour),
		Duration:  30 * time.Minute,
		UserID:    userIDTwo,
	}
	event3 := storage.Event{
		Title:       "Event 3",
		StartTime:   dto.now.Add(50 * time.Minute),
		Duration:    30 * time.Minute,
		Description: "Test event 3",
		UserID:      userIDOne,
	}
	event4 := storage.Event{
		StartTime:   dto.now.Add(-50 * time.Minute),
		Duration:    30 * time.Minute,
		Description: "Test event 4",
	}
	event5 := storage.Event{
		ID:          wrongEventId,
		StartTime:   dto.now.Add(-50 * time.Minute),
		Duration:    30 * time.Minute,
		Description: "Test event 5",
	}
	dto.events = append(dto.events, event0, event1, event2, event3, event4, event5)
	return dto
}
