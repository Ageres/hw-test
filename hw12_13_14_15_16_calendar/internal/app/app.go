package app

import (
	"context"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	storage storage.Storage
}

type Status string

const (
	ADD       Status = "Event added successfully"
	UPDATE    Status = "Event updated successfully"
	DELETE    Status = "Event deleted successfully"
	LISTDAY   Status = "Day event list successfully retrieved"
	LISTWEEK  Status = "Week event list successfully retrieved"
	LISTMONTH Status = "Month event list successfully retrieved"
)

type Response struct {
	Status Status `json:"status" binding:"required"`
	Result any    `json:"result,omitempty"`
}

func New(_ context.Context, storage storage.Storage) *App {
	return &App{storage: storage}
}

func (a *App) AddEvent(ctx context.Context, eventRef *storage.Event) (*Response, error) {
	respEventRef, err := a.storage.Add(ctx, eventRef)
	if err != nil {
		return nil, err
	}
	return &Response{
		Status: ADD,
		Result: respEventRef,
	}, err
}

func (a *App) UpdateEvent(ctx context.Context, eventRef *storage.Event) (*Response, error) {
	err := a.storage.Update(ctx, eventRef)
	if err != nil {
		return nil, err
	}
	return &Response{
		Status: UPDATE,
	}, err
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) (*Response, error) {
	err := a.storage.Delete(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return &Response{
		Status: DELETE,
	}, err
}

func (a *App) ListDayEvents(ctx context.Context, startDay time.Time) (*Response, error) {
	respEvents, err := a.storage.ListDay(ctx, startDay)
	if err != nil {
		return nil, err
	}
	return &Response{
		Status: LISTDAY,
		Result: respEvents,
	}, err
}

func (a *App) ListWeekEvents(ctx context.Context, startDay time.Time) (*Response, error) {
	respEvents, err := a.storage.ListWeek(ctx, startDay)
	if err != nil {
		return nil, err
	}
	return &Response{
		Status: LISTWEEK,
		Result: respEvents,
	}, err
}

func (a *App) ListMonthEvents(ctx context.Context, startDay time.Time) (*Response, error) {
	respEvents, err := a.storage.ListMonth(ctx, startDay)
	if err != nil {
		return nil, err
	}
	return &Response{
		Status: LISTMONTH,
		Result: respEvents,
	}, err
}
