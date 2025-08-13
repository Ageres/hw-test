package storage

import (
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
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

func (e *Event) ToNotification() *model.Notification {
	return &model.Notification{
		ID:        e.ID,
		Title:     e.Title,
		StartTime: e.StartTime,
		UserID:    e.UserID,
	}
}
