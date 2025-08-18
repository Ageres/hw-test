package app

import (
	"context"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(_ context.Context, _ Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(_ context.Context, id, title string) error {
	// TODO
	_ = id
	_ = title

	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
