package storage

import (
	"errors"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

var (
	ErrDateBusy      = errors.New("time is already taken by another event")
	ErrNotifyTooLate = errors.New("notification time has already expired")
	ErrUserConflict  = errors.New("user is not the owner of the event")
	ErrEventNotFound = errors.New("no events found")
)

type Event struct {
	ID          string //UUID
	Title       string
	StartTime   time.Time
	Duration    time.Duration
	Description string
	UserID      string //UUID
	Reminder    time.Duration
}

type Storage interface {
	Add(event Event) error
	Update(id string, event Event) error
	Delete(id string) error
	ListDay(day time.Time) ([]Event, error)
	ListWeek(start time.Time) ([]Event, error)
	ListMonth(start time.Time) ([]Event, error)
}

func (e *Event) ToNotification() *model.Notification {
	return &model.Notification{
		ID:        e.ID,
		Title:     e.Title,
		StartTime: e.StartTime,
		UserID:    e.UserID,
	}
}
