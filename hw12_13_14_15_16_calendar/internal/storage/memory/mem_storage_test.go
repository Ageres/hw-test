package memorystorage

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
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
		require.NoError(t, dto.storage.Add(&events[0]))
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[6]))

		weekEvents, err := dto.storage.ListWeek(dto.now)
		require.NoError(t, err)
		require.Len(t, weekEvents, 2)

		for _, event := range weekEvents {
			require.NotEqual(t, "Next Week Event", event.Title)
		}
	})

	t.Run("event spans week boundary", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[7]))

		weekEvents, err := dto.storage.ListWeek(dto.now)
		require.NoError(t, err)
		require.Len(t, weekEvents, 1)
	})

	t.Run("list month events", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[0]))
		require.NoError(t, dto.storage.Add(&events[1]))
		require.NoError(t, dto.storage.Add(&events[8]))

		monthEvents, err := dto.storage.ListMonth(dto.now)
		require.NoError(t, err)
		require.Len(t, monthEvents, 2)
	})

	t.Run("event spans month boundary", func(t *testing.T) {
		dto.buildNewStorage()
		require.NoError(t, dto.storage.Add(&events[9]))

		monthEvents, err := dto.storage.ListMonth(dto.now)
		require.NoError(t, err)
		require.Len(t, monthEvents, 1)
	})
}

func TestStorageConcurrentAdd(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewStorage()
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 1; j <= 10; j++ {
				event := storage.Event{
					Title:     fmt.Sprintf("Event %d-%d", idx, j),
					StartTime: time.Now().Add(time.Duration(j) * time.Hour),
					Duration:  30 * time.Minute,
					UserID:    uuid.New().String(),
				}
				err := dto.storage.Add(&event)
				require.NoError(t, err)
			}
		}(i)
	}
	wg.Wait()

	require.Len(t, dto.storage.events, 100)
}

func TestStorageConcurrentReadWrite(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewStorage()
	dto.now = time.Now()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 100 {
			event1 := storage.Event{
				Title:     "Event 1",
				StartTime: dto.now.Add(2 * time.Hour),
				Duration:  1 * time.Hour,
				UserID:    strconv.Itoa(i),
			}
			err := dto.storage.Add(&event1)
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 20 {
				_, err := dto.storage.ListDay(time.Now())
				require.NoError(t, err)
				time.Sleep(15 * time.Millisecond)
			}
		}()
	}
	wg.Wait()
}

func TestStorageConcurrentUpdateDelete(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewStorage()
	events := dto.buildNewEvents().events
	require.NoError(t, dto.storage.Add(&events[0]))
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := range 10 {
				updatedEvent := events[0]
				updatedEvent.Title = fmt.Sprintf("Updated %d-%d", idx, j)
				err := dto.storage.Update(events[0].ID, &updatedEvent)
				require.NoError(t, err)
			}
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		err := dto.storage.Delete(events[0].ID)
		require.NoError(t, err)
	}()

	wg.Wait()
	_, exists := dto.storage.events[events[0].ID]
	require.False(t, exists)
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
	event6 := storage.Event{
		Title:     "Event 6 - Next Week Event",
		StartTime: dto.now.Add(7 * 24 * time.Hour), // +7 дней
		Duration:  1 * time.Hour,
		UserID:    userIDOne,
	}
	endOfWeek := dto.now.Add(time.Duration(6-int(dto.now.Weekday())) * 24 * time.Hour)
	event7 := storage.Event{
		Title:     "event7 - Week Boundary Event",
		StartTime: endOfWeek.Add(-12 * time.Hour),
		Duration:  36 * time.Hour,
		UserID:    userIDOne,
	}
	nextMonth := dto.now.AddDate(0, 1, 0)
	event8 := storage.Event{
		Title:     "event8 - Next Month Event",
		StartTime: nextMonth,
		Duration:  1 * time.Hour,
		UserID:    userIDOne,
	}
	endOfMonth := time.Date(dto.now.Year(), dto.now.Month()+1, 0, 0, 0, 0, 0, dto.now.Location())
	event9 := storage.Event{
		Title:     "event9 - Month Boundary Event",
		StartTime: endOfMonth.Add(-12 * time.Hour),
		Duration:  36 * time.Hour,
		UserID:    userIDOne,
	}
	dto.events = append(dto.events, event0, event1, event2, event3, event4, event5, event6, event7, event8, event9)
	return dto
}
