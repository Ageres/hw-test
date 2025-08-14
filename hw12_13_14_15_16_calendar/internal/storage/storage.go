package storage

import (
	"context"
	"time"
)

type Storage interface {
	Add(ctx context.Context, eventRef *Event) error
	Update(ctx context.Context, eventRef *Event) error
	Delete(ctx context.Context, id string) error
	ListDay(ctx context.Context, startDay time.Time) ([]Event, error)
	ListWeek(ctx context.Context, startDay time.Time) ([]Event, error)
	ListMonth(ctx context.Context, startDay time.Time) ([]Event, error)
}
