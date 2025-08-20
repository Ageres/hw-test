package app

import (
	"context"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	storage storage.Storage
}

type AppResponse struct {
	Status string
	Result any
}

func New(_ context.Context, storage storage.Storage) *App {
	return &App{storage: storage}
}

func (a *App) AddEvent(ctx context.Context, eventRef *storage.Event) (*storage.Event, error) {
	return a.storage.Add(ctx, eventRef)
}

func (a *App) UpdateEvent(ctx context.Context, eventRef *storage.Event) error {
	return a.storage.Update(ctx, eventRef)
}

func (a *App) DeleteEvent(ctx context.Context, eventId string) error {
	return a.storage.Delete(ctx, eventId)
}

func (a *App) ListDayEvents(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	return a.storage.ListDay(ctx, startDay)
}

func (a *App) ListWeekEvents(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	return a.storage.ListWeek(ctx, startDay)
}

func (a *App) ListMonthEvents(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	return a.storage.ListMonth(ctx, startDay)
}
