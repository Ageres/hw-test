package utils

import (
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

func EventToNotification(e *storage.Event) *model.Notification {
	if e == nil {
		return nil
	}
	return &model.Notification{
		ID:    e.ID,
		Title: e.Title,
		Date: time.Date(
			e.StartTime.Year(),
			e.StartTime.Month(),
			e.StartTime.Day(),
			0, 0, 0, 0,
			e.StartTime.Location(),
		),
		UserID: e.UserID,
	}
}
