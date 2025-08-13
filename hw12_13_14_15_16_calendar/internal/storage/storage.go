package storage

import (
	"errors"
	"time"
)

var (
	ErrDateBusy      = errors.New("time is already taken by another event")
	ErrNotifyTooLate = errors.New("notification time has already expired")
	ErrUserConflict  = errors.New("user is not the owner of the event")
	ErrEventNotFound = errors.New("no events found")
)

type Storage interface {
	Get(id string) (Event, error)
	Add(event Event) error
	Update(id string, event Event) error
	Delete(id string) error
	ListDay(day time.Time) ([]Event, error)
	ListWeek(start time.Time) ([]Event, error)
	ListMonth(start time.Time) ([]Event, error)
	ListByUser(userId string) ([]Event, error)
}
