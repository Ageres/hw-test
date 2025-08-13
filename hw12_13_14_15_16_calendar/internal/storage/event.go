package storage

import "time"

type Event struct {
	ID          string //UUID
	Title       string
	StartTime   time.Time
	Duration    time.Duration
	Description string
	UserID      string //UUID
	Reminder    time.Duration
}

func (e *Event) ToNotification() Notification {
	return Notification{
		ID:        e.ID,
		Title:     e.Title,
		StartTime: e.StartTime,
		UserID:    e.UserID,
	}
}
