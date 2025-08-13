package storage

import (
	"time"
)

type Storage interface {
	Add(eventRef *Event) error
	Update(id string, eventRef *Event) error
	Delete(id string) error
	ListDay(start time.Time) ([]Event, error)
	ListWeek(start time.Time) ([]Event, error)
	ListMonth(start time.Time) ([]Event, error)
}
