package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	/*
		"fmt"
		"strconv"
		"sync"
		"testing"
		"time"
		"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
		"github.com/google/uuid"
		"github.com/stretchr/testify/require"
	*/)

func TestStorageAdd(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("add event", func(t *testing.T) {
		dto.buildNewStorage()
		require.Len(t, dto.storage.events, 0)

		oldGenerator := storage.FnUuidGenerator
		defer func() { storage.FnUuidGenerator = oldGenerator }()

		mockUUIDs := []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		}
		callCount := 0
		storage.FnUuidGenerator = func() uuid.UUID {
			if callCount >= len(mockUUIDs) {
				t.Fatal("Unexpected call to UUID generator")
			}
			u := mockUUIDs[callCount]
			callCount++
			return u
		}

		testEvents := []storage.Event{events[0], events[1], events[2]}
		for i := range testEvents {
			testEvents[i].ID = ""
		}

		var err error
		addedEvents := make([]*storage.Event, 3)
		for i := range testEvents {
			addedEvents[i], err = dto.storage.Add(dto.testContext, &testEvents[i])
			require.NoError(t, err)
		}

		require.Len(t, dto.storage.events, 3)
		require.Equal(t, 3, callCount, "UUID generator should be called 3 times")

		for i, event := range addedEvents {

			expectedID := mockUUIDs[i].String()
			require.Equal(t, expectedID, event.ID, "Event ID should match mock UUID")

			storedEvent, exists := dto.storage.events[expectedID]
			require.True(t, exists, "Event should exist in storage")

			require.Equal(t, testEvents[i].Title, storedEvent.Title)
			require.Equal(t, testEvents[i].StartTime, storedEvent.StartTime)
			require.Equal(t, testEvents[i].Duration, storedEvent.Duration)
			require.Equal(t, testEvents[i].Description, storedEvent.Description)
			require.Equal(t, testEvents[i].UserID, storedEvent.UserID)
			require.Equal(t, testEvents[i].Reminder, storedEvent.Reminder)

			require.Equal(t, storedEvent, *addedEvents[i])
		}
	})

	t.Run("nil event error when adding", func(t *testing.T) {
		dto.buildNewStorage()
		_, err := dto.storage.Add(dto.testContext, nil)
		require.ErrorIs(t, err, storage.ErrEventIsNil)
		require.Len(t, dto.storage.events, 0)
	})

	t.Run("validation event error when adding", func(t *testing.T) {
		dto.buildNewStorage()
		_, err := dto.storage.Add(dto.testContext, &events[0])
		require.NoError(t, err)

		_, err = dto.storage.Add(dto.testContext, &events[4])
		require.Error(t, err)
		require.Contains(t, err.Error(), "title is empty")
		require.Contains(t, err.Error(), "event time is expired")
		require.Contains(t, err.Error(), "user id is empty")
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("event duplication error when adding", func(t *testing.T) {
		dto.buildNewStorage()

		oldGenerator := storage.FnUuidGenerator
		defer func() { storage.FnUuidGenerator = oldGenerator }()

		fixedUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		storage.FnUuidGenerator = func() uuid.UUID { return fixedUUID }

		testEvent := storage.Event{
			Title:     "Test Event",
			StartTime: time.Now().Add(1 * time.Hour),
			Duration:  30 * time.Minute,
			UserID:    "test-user",
		}

		addedEvent, err := dto.storage.Add(dto.testContext, &testEvent)
		require.NoError(t, err)
		require.Equal(t, fixedUUID.String(), addedEvent.ID)
		require.Len(t, dto.storage.events, 1)

		_, err = dto.storage.Add(dto.testContext, &testEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "event with this id already exists")
		require.Len(t, dto.storage.events, 1)

		storedEvent, exists := dto.storage.events[fixedUUID.String()]
		require.True(t, exists)
		require.Equal(t, addedEvent, &storedEvent)
	})

	t.Run("date busy error when adding", func(t *testing.T) {
		dto.buildNewStorage()
		_, err := dto.storage.Add(dto.testContext, &events[0])
		require.NoError(t, err)

		_, err = dto.storage.Add(dto.testContext, &events[3])
		require.Error(t, err)
		require.Contains(t, err.Error(), "time is already taken")
		require.Len(t, dto.storage.events, 1)
	})

	t.Run("generated ID is valid UUID", func(t *testing.T) {
		dto.buildNewStorage()

		oldGenerator := storage.FnUuidGenerator
		defer func() { storage.FnUuidGenerator = oldGenerator }()

		fixedUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		storage.FnUuidGenerator = func() uuid.UUID { return fixedUUID }

		event := storage.Event{
			Title:     "Test Event",
			StartTime: time.Now().Add(1 * time.Hour),
			Duration:  30 * time.Minute,
			UserID:    "test-user",
		}

		result, err := dto.storage.Add(dto.testContext, &event)
		require.NoError(t, err)
		require.Equal(t, fixedUUID.String(), result.ID)
	})
}

func TestStorageUpdate(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("update event", func(t *testing.T) {
		dto.buildNewStorage()

		oldGenerator := storage.FnUuidGenerator
		defer func() { storage.FnUuidGenerator = oldGenerator }()

		fixedUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		storage.FnUuidGenerator = func() uuid.UUID { return fixedUUID }

		originalEvent := storage.Event{
			Title:       "Original Event",
			StartTime:   time.Now().Add(1 * time.Hour),
			Duration:    30 * time.Minute,
			Description: "Original Description",
			UserID:      "user-1",
			Reminder:    15 * time.Minute,
		}

		_, err := dto.storage.Add(dto.testContext, &originalEvent)
		require.NoError(t, err)
		require.Len(t, dto.storage.events, 1)
		require.Equal(t, fixedUUID.String(), originalEvent.ID)

		updatedEvent := storage.Event{
			ID:          fixedUUID.String(),
			Title:       "Updated Event",
			StartTime:   time.Now().Add(2 * time.Hour),
			Duration:    1 * time.Hour,
			Description: "Updated Description",
			UserID:      "user-1",
			Reminder:    30 * time.Minute,
		}

		err = dto.storage.Update(dto.testContext, &updatedEvent)
		require.NoError(t, err)
		require.Len(t, dto.storage.events, 1)

		storedEvent, exists := dto.storage.events[fixedUUID.String()]
		require.True(t, exists)

		require.Equal(t, updatedEvent.ID, storedEvent.ID, "ID should remain the same")
		require.Equal(t, updatedEvent.Title, storedEvent.Title, "Title should be updated")
		require.Equal(t, updatedEvent.StartTime, storedEvent.StartTime, "StartTime should be updated")
		require.Equal(t, updatedEvent.Duration, storedEvent.Duration, "Duration should be updated")
		require.Equal(t, updatedEvent.Description, storedEvent.Description, "Description should be updated")
		require.Equal(t, updatedEvent.UserID, storedEvent.UserID, "UserID should remain the same")
		require.Equal(t, updatedEvent.Reminder, storedEvent.Reminder, "Reminder should be updated")

		require.NotEqual(t, originalEvent.Title, storedEvent.Title)
		require.NotEqual(t, originalEvent.StartTime, storedEvent.StartTime)
		require.NotEqual(t, originalEvent.Duration, storedEvent.Duration)
		require.NotEqual(t, originalEvent.Description, storedEvent.Description)
		require.NotEqual(t, originalEvent.Reminder, storedEvent.Reminder)
	})

	t.Run("nil event error when updating", func(t *testing.T) {
		dto.buildNewStorage()
		_, err := dto.storage.Add(dto.testContext, &events[0])
		require.NoError(t, err)

		err = dto.storage.Update(dto.testContext, nil)
		require.ErrorIs(t, err, storage.ErrEventIsNil)
	})

	t.Run("validation event error when updating", func(t *testing.T) {
		dto.buildNewStorage()
		_, err := dto.storage.Add(dto.testContext, &events[0])
		require.NoError(t, err)

		invalidEvent := events[4]
		invalidEvent.ID = events[0].ID
		err = dto.storage.Update(dto.testContext, &invalidEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "title is empty")
	})

	t.Run("event not found error when updating", func(t *testing.T) {
		dto.buildNewStorage()
		err := dto.storage.Update(dto.testContext, &events[0])
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("user conflict error when updating", func(t *testing.T) {
		dto.buildNewStorage()
		_, err := dto.storage.Add(dto.testContext, &events[0])
		require.NoError(t, err)

		conflictEvent := events[2]
		conflictEvent.ID = events[0].ID
		err = dto.storage.Update(dto.testContext, &conflictEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "is not the owner of the event")
	})

	t.Run("date busy error when updating", func(t *testing.T) {
		dto.buildNewStorage()

		oldGenerator := storage.FnUuidGenerator
		defer func() { storage.FnUuidGenerator = oldGenerator }()

		fixedUUIDs := []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		}
		callCount := 0
		storage.FnUuidGenerator = func() uuid.UUID {
			u := fixedUUIDs[callCount]
			callCount++
			return u
		}

		event1 := storage.Event{
			Title:     "Event 1",
			StartTime: time.Now().Add(1 * time.Hour),
			Duration:  30 * time.Minute,
			UserID:    "user-1",
		}

		event2 := storage.Event{
			Title:     "Event 2",
			StartTime: time.Now().Add(2 * time.Hour),
			Duration:  30 * time.Minute,
			UserID:    "user-1",
		}

		_, err := dto.storage.Add(dto.testContext, &event1)
		require.NoError(t, err)
		_, err = dto.storage.Add(dto.testContext, &event2)
		require.NoError(t, err)
		require.Len(t, dto.storage.events, 2)

		updatedEvent2 := event2
		updatedEvent2.StartTime = event1.StartTime.Add(15 * time.Minute)

		err = dto.storage.Update(dto.testContext, &updatedEvent2)
		require.Error(t, err)
		require.Equal(t, err.Error(), "time is already taken by another event: 11111111-1111-1111-1111-111111111111")

		storedEvent1, exists1 := dto.storage.events[fixedUUIDs[0].String()]
		require.True(t, exists1)
		require.Equal(t, event1, storedEvent1)

		storedEvent2, exists2 := dto.storage.events[fixedUUIDs[1].String()]
		require.True(t, exists2)
		require.Equal(t, event2, storedEvent2)
	})

}

func TestStorageDelete(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	events := dto.events

	t.Run("delete event", func(t *testing.T) {
		dto.buildNewStorage()

		oldGenerator := storage.FnUuidGenerator
		defer func() { storage.FnUuidGenerator = oldGenerator }()

		fixedUUIDs := []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		}
		callCount := 0
		storage.FnUuidGenerator = func() uuid.UUID {
			u := fixedUUIDs[callCount]
			callCount++
			return u
		}

		event1 := storage.Event{
			Title:     "Event 1",
			StartTime: time.Now().Add(1 * time.Hour),
			Duration:  30 * time.Minute,
			UserID:    "user-1",
		}

		event2 := storage.Event{
			Title:     "Event 2",
			StartTime: time.Now().Add(2 * time.Hour),
			Duration:  1 * time.Hour,
			UserID:    "user-1",
		}

		addedEvent1, err := dto.storage.Add(dto.testContext, &event1)
		require.NoError(t, err)
		addedEvent2, err := dto.storage.Add(dto.testContext, &event2)
		require.NoError(t, err)
		require.Len(t, dto.storage.events, 2)

		err = dto.storage.Delete(dto.testContext, addedEvent1.ID)
		require.NoError(t, err)
		require.Len(t, dto.storage.events, 1)

		remainingEvent, exists := dto.storage.events[addedEvent2.ID]
		require.True(t, exists, "Second event should remain in storage")

		require.Equal(t, addedEvent2.ID, remainingEvent.ID, "ID should match")
		require.Equal(t, event2.Title, remainingEvent.Title, "Title should match")
		require.Equal(t, event2.StartTime, remainingEvent.StartTime, "StartTime should match")
		require.Equal(t, event2.Duration, remainingEvent.Duration, "Duration should match")
		require.Equal(t, event2.Description, remainingEvent.Description, "Description should match")
		require.Equal(t, event2.UserID, remainingEvent.UserID, "UserID should match")
		require.Equal(t, event2.Reminder, remainingEvent.Reminder, "Reminder should match")

		_, exists = dto.storage.events[addedEvent1.ID]
		require.False(t, exists, "First event should be deleted")
	})

	t.Run("validation id error when deleting", func(t *testing.T) {
		dto.buildNewStorage()
		_, err := dto.storage.Add(dto.testContext, &events[0])
		require.NoError(t, err)

		err = dto.storage.Delete(dto.testContext, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "validate event id")
	})

	t.Run("event not found error when deleting", func(t *testing.T) {
		dto.buildNewStorage()
		err := dto.storage.Delete(dto.testContext, events[0].ID)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}

func TestStorageListEvents(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewEvents()
	//events := dto.events

	t.Run("list day events", func(t *testing.T) {
		dto.buildNewStorage()

		oldGenerator := storage.FnUuidGenerator
		defer func() { storage.FnUuidGenerator = oldGenerator }()

		fixedUUIDs := []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		}
		callCount := 0
		storage.FnUuidGenerator = func() uuid.UUID {
			u := fixedUUIDs[callCount]
			callCount++
			return u
		}

		startDate := time.Now().Add(1 * time.Hour)
		event1 := storage.Event{
			Title:     "Morning Event",
			StartTime: time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 10, 0, 0, 0, startDate.Location()),
			Duration:  1 * time.Hour,
			UserID:    "user-1",
		}

		event2 := storage.Event{
			Title:     "Evening Event",
			StartTime: time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 18, 0, 0, 0, startDate.Location()),
			Duration:  2 * time.Hour,
			UserID:    "user-1",
		}

		addedEvent1, err := dto.storage.Add(dto.testContext, &event1)
		require.NoError(t, err)
		addedEvent2, err := dto.storage.Add(dto.testContext, &event2)
		require.NoError(t, err)

		dayEvents, err := dto.storage.ListDay(dto.testContext, startDate)
		require.NoError(t, err)
		require.Len(t, dayEvents, 2)

		for _, returnedEvent := range dayEvents {
			switch returnedEvent.ID {
			case addedEvent1.ID:
				require.Equal(t, event1.Title, returnedEvent.Title)
				require.Equal(t, event1.StartTime, returnedEvent.StartTime)
				require.Equal(t, event1.Duration, returnedEvent.Duration)
				require.Equal(t, event1.UserID, returnedEvent.UserID)
			case addedEvent2.ID:
				require.Equal(t, event2.Title, returnedEvent.Title)
				require.Equal(t, event2.StartTime, returnedEvent.StartTime)
				require.Equal(t, event2.Duration, returnedEvent.Duration)
				require.Equal(t, event2.UserID, returnedEvent.UserID)
			default:
				t.Fatalf("Unexpected event ID: %s", returnedEvent.ID)
			}
		}
	})

	/*
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
	*/
}

/*
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

func TestStorageDeadlock(t *testing.T) {
	dto := newTestMemoryStorageDto().buildNewStorage()
	event := dto.buildNewEvents().events[0]
	require.NoError(t, dto.storage.Add(&event))

	done := make(chan bool)
	go func() {
		dto.storage.Update(event.ID, &event)
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("Potential deadlock detected")
	}
}
*/

func TestGenerateTestEvents(t *testing.T) {
	storage := &MemoryStorage{
		events: make(map[string]storage.Event),
	}

	/*
		startTime := time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC)
		period := time.Hour
		userCount := 10
		eventsPerUser := 5
	*/

	storage.generateTestEvents()

	count := 0
	for _, event := range storage.events {
		/*
			if count >= 3 {
				break
			}
		*/
		fmt.Printf("Event: %+v\n", event)
		count++
	}

}

// -------------------------------------------------------------------------------------
// Вспомогательные функции
type TestMemoryStorageDto struct {
	storage     *MemoryStorage
	events      []storage.Event
	now         time.Time
	testContext context.Context
}

func newTestMemoryStorageDto() *TestMemoryStorageDto {
	ctx := logger.SetDefaultLogger(context.Background())
	return &TestMemoryStorageDto{
		testContext: ctx,
		now:         time.Now(),
	}
}

func (dto *TestMemoryStorageDto) buildNewStorage() *TestMemoryStorageDto {
	storageConfRef := model.StorageConf{
		LoadTestData: false,
	}
	dto.storage = NewMemoryStorage(dto.testContext, &storageConfRef).(*MemoryStorage)
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
