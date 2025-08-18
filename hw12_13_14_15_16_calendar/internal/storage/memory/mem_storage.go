package memorystorage

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	events map[string]storage.Event // key: Event.ID
}

func NewMemoryStorage(ctx context.Context, storageConfRef *model.StorageConf) storage.Storage {
	storage := &MemoryStorage{
		events: make(map[string]storage.Event),
	}
	if storageConfRef.LoadTestData {
		storage.generateTestEvents()
		lg.GetLogger(ctx).Info("test event loaded")
	}
	return storage
}

func (m *MemoryStorage) Add(ctx context.Context, eventRef *storage.Event) (*storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("add event", map[string]any{"event": lg.MarshalAny(eventRef)})

	if err := storage.ValidateEvent(eventRef); err != nil {
		logger.WithError(err).Error("add event")
		return nil, err
	}
	eventRef.GenerateEventID()

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.events[eventRef.ID]; exists {
		err := storage.NewSErrorWithTemplate("failed to add event, event with this id already exists: %s", eventRef.ID)
		logger.WithError(err).Error("add event")
		return nil, err
	}

	for _, existingEvent := range m.events {
		if existingEvent.UserID == eventRef.UserID &&
			overlaps(&existingEvent, eventRef) {
			err := storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, existingEvent.ID)
			logger.WithError(err).Error("add event")
			return nil, err
		}
	}

	m.events[eventRef.ID] = *eventRef
	result := *eventRef
	return &result, nil
}

func (m *MemoryStorage) Update(ctx context.Context, eventRef *storage.Event) error {
	logger := lg.GetLogger(ctx)
	logger.Info("update event", map[string]any{"event": lg.MarshalAny(eventRef)})

	if err := storage.FullValidateEvent(eventRef); err != nil {
		logger.WithError(err).Error("update event")
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id := eventRef.ID
	oldEvent, exists := m.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}
	if oldEvent.UserID != eventRef.UserID {
		err := storage.NewSErrorWithTemplate(storage.ErrUserConflictMsgTemplate, eventRef.UserID, oldEvent.UserID)
		logger.WithError(err).Error("update event")
		return err
	}

	for _, existingEvent := range m.events {
		if existingEvent.UserID == eventRef.UserID &&
			existingEvent.ID != id &&
			overlaps(&existingEvent, eventRef) {
			err := storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, existingEvent.ID)
			logger.WithError(err).Error("update event")
			return err
		}
	}

	m.events[id] = *eventRef
	return nil
}

func (m *MemoryStorage) Delete(ctx context.Context, id string) error {
	logger := lg.GetLogger(ctx)
	logger.Info("delete event", map[string]any{"eventId": id})

	if err := storage.ValidateEventID(id); err != nil {
		logger.WithError(err).Error("delete event")
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}

	delete(m.events, id)
	return nil
}

func (m *MemoryStorage) ListDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list day events", map[string]any{"startDay": startDay})

	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.Add(24 * time.Hour)
	return m.listEvents(ctx, startTime, endTime)
}

func (m *MemoryStorage) ListWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list week events", map[string]any{"startDay": startDay})

	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.AddDate(0, 0, 7)
	return m.listEvents(ctx, startTime, endTime)
}

func (m *MemoryStorage) ListMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list month events", map[string]any{"startDay": startDay})

	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.AddDate(0, 1, 0)
	return m.listEvents(ctx, startTime, endTime)
}

func (m *MemoryStorage) listEvents(ctx context.Context, startTime, endTime time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list events", map[string]any{
		"startTime": startTime,
		"endTime":   endTime,
	})

	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []storage.Event
	for _, event := range m.events {
		select {
		case <-ctx.Done():
			err := storage.NewSError(storage.ErrContextDone, ctx.Err())
			logger.WithError(err).Error("list events")
			return nil, err
		default:
		}

		eventEnd := event.StartTime.Add(event.Duration)
		if (event.StartTime.After(startTime) && event.StartTime.Before(endTime)) ||
			(eventEnd.After(startTime) && eventEnd.Before(endTime)) {
			result = append(result, event)
		}
	}
	logger.Info("list events", map[string]any{"found": len(result)})
	return result, nil
}

func (m *MemoryStorage) generateTestEvents() {
	// параметры генерации
	startTime := time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC)
	period := 12 * time.Hour
	userCount := 100
	eventsPerUser := 100

	eventID := 0
	for userID := 1; userID <= userCount; userID++ {
		userIDStr := fmt.Sprintf("user-%04d", userID)
		currentTime := startTime

		for eventNum := 1; eventNum <= eventsPerUser; eventNum++ {
			// генерация ID эвента
			eventIDStr := fmt.Sprintf("%08d-%04d-%04d-%04d-%012d",
				0, 0, 0, 0, eventID)
			eventID++
			// генерация тайтла эвента
			title := fmt.Sprintf("title_%s_%d", userIDStr, eventNum)
			// генерация описания эвента
			description := fmt.Sprintf("%s_desc", title)
			// генерация случайной длительности эвента (1 мин - 2 суток)
			duration := time.Duration(rand.Int63n(2*24*60*60-60)+1) * time.Second
			// генерация случайного периода напоминания до эвента (1 мин - 2 суток)
			reminder := time.Duration(rand.Int63n(2*24*60*60-60)+1) * time.Second
			// эвент
			event := storage.Event{
				ID:          eventIDStr,
				Title:       title,
				StartTime:   currentTime,
				Duration:    duration,
				Description: description,
				UserID:      userIDStr,
				Reminder:    reminder,
			}
			m.events[eventIDStr] = event
			// время для следующего события
			currentTime = currentTime.Add(period)
		}
	}
}

// проверка двух эвентов на пересечение времени.
func overlaps(e, other *storage.Event) bool {
	if e == nil || other == nil {
		return false
	}
	end1 := e.StartTime.Add(e.Duration)
	end2 := other.StartTime.Add(other.Duration)
	return e.StartTime.Before(end2) && end1.After(other.StartTime)
}
