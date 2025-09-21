package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	mu        sync.RWMutex
	events    map[string]storage.Event // key: Event.ID
	procEvent []string
}

func NewMemoryStorage(ctx context.Context, storageConfRef *model.StorageConf) storage.Storage {
	storage := &MemoryStorage{
		events: make(map[string]storage.Event),
	}
	if storageConfRef.InMemory.LoadTestData {
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
	userCount := 100
	eventsPerUser := 100
	eventCounter := 0

	for userID := 1; userID <= userCount; userID++ {
		userIDStr := fmt.Sprintf("user-%04d", userID)

		for eventNum := 1; eventNum <= eventsPerUser; eventNum++ {
			eventID := fmt.Sprintf("00000000-0000-0000-0000-%012d", eventCounter)
			eventCounter++

			title := fmt.Sprintf("title_%s_%d", userIDStr, eventNum)
			description := fmt.Sprintf("%s_desc", title)

			duration := time.Duration(rand.Int63n(2*24*60*60-60)+1) * time.Second //nolint:gosec // генерация тестовых данных

			startOffset := rand.Float64()*1095 - 547.5 //nolint:gosec // генерация тестовых данных
			startTime := time.Now().Add(time.Duration(startOffset * 24 * float64(time.Hour)))

			reminder := time.Duration(rand.Int63n(2*24*60*60-60)+1) * time.Second //nolint:gosec // генерация тестовых данных

			event := storage.Event{
				ID:          eventID,
				Title:       title,
				StartTime:   startTime,
				Duration:    duration,
				Description: description,
				UserID:      userIDStr,
				Reminder:    reminder,
			}

			m.events[eventID] = event
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

func (m *MemoryStorage) ListReminderEvents(ctx context.Context, scanInterval int64) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list reminder events", map[string]any{
		"scanInterval": scanInterval,
	})

	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []storage.Event
	for _, event := range m.events {
		select {
		case <-ctx.Done():
			err := storage.NewSError(storage.ErrContextDone, ctx.Err())
			logger.WithError(err).Error("list reminder events")
			return nil, err
		default:
		}

		now := time.Now()
		if event.Reminder == 0 || event.StartTime.Before(now) {
			continue
		}
		reminderTime := event.StartTime.Add(-event.Reminder)
		scanTime := now.Add(time.Duration(scanInterval) * time.Second)
		if reminderTime.Before(scanTime) {
			result = append(result, event)
		}
	}
	logger.Info("list reminder events", map[string]any{"found": len(result)})
	return result, nil
}

func (m *MemoryStorage) ResetEventReminder(ctx context.Context, eventIDs []string) error {
	logger := lg.GetLogger(ctx)

	eventIDsLen := len(eventIDs)
	logger.Info("reset event reminder", map[string]any{"eventIDLen": eventIDsLen})

	if eventIDsLen == 0 {
		err := storage.ErrEventIDListIsEmpty
		logger.WithError(err).Error("reset event reminder")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range eventIDs {
		if err := storage.ValidateEventID(id); err != nil {
			logger.WithError(err).Error("reset event reminder")
			return err
		}
		event, exists := m.events[id]
		if !exists {
			return storage.ErrEventNotFound
		}
		event.Reminder = 0
		m.events[id] = event
	}

	return nil
}

func (m *MemoryStorage) DeleteOldEvents(ctx context.Context, before time.Time) (int64, error) {
	logger := lg.GetLogger(ctx)

	logger.Info("delete old events", map[string]any{"before": before})

	m.mu.Lock()
	defer m.mu.Unlock()
	var rows int64
	for _, event := range m.events {
		if event.StartTime.Before(before) {
			delete(m.events, event.ID)
			rows++
		}
	}

	if rows == 0 {
		logger.WithError(storage.ErrEventNotFound).Debug("delete old events")
	}

	return rows, nil
}

func (m *MemoryStorage) AddProcEvent(ctx context.Context, procEventRef *storage.ProcEvent) error {
	logger := lg.GetLogger(ctx)
	logger.Info("add proc event", map[string]any{"procEvent": lg.MarshalAny(procEventRef)})

	if procEventRef == nil {
		err := errors.New("procEvent is nil")
		logger.WithError(err).Error("add proc event", map[string]any{"procEvent": procEventRef})
		return err
	}

	procEventID := procEventRef.ID

	err := uuid.Validate(procEventID)
	if err != nil {
		logger.WithError(err).Error("add proc event", map[string]any{"procEvent": procEventRef})
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if slices.Contains(m.procEvent, procEventID) {
		err := storage.NewSErrorWithTemplate("failed to add proc event, proc event with this id already exists: %s", procEventID)
		logger.WithError(err).Error("add event")
		return err
	}

	m.procEvent = append(m.procEvent, procEventID)
	return nil
}

func (m *MemoryStorage) Close() error {
	return nil
}
