package storage

import (
	"fmt"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

var FnUUIDGenerator = uuid.New

type Event struct {
	ID          string // UUID
	Title       string
	StartTime   time.Time
	Duration    time.Duration
	Description string
	UserID      string
	Reminder    time.Duration
}

func (e *Event) ToNotification() *model.Notification {
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

func (e *Event) GenerateEventId() {
	e.ID = FnUUIDGenerator().String()
}

func ValidateEventID(eventID string) error {
	err := uuid.Validate(eventID)
	if err != nil {
		return NewSError("failed to validate event id", err)
	}
	return nil
}

// без валидации ID.
func ValidateEvent(e *Event) error {
	if e == nil {
		return ErrEventIsNil
	}
	errMsgs := e.simpleValidate()
	return NewSErrorWithMsgArr(errMsgs)
}

func FullValidateEvent(e *Event) error {
	if e == nil {
		return ErrEventIsNil
	}
	errMsgs := make([]string, 0, 5)
	err := uuid.Validate(e.ID)
	if err != nil {
		errMsgs = append(errMsgs, fmt.Sprintf("failed to validate event id: %v", err))
	}
	errMsgs = append(errMsgs, e.simpleValidate()...)
	return NewSErrorWithMsgArr(errMsgs)
}

func (e *Event) simpleValidate() []string {
	errMsgs := make([]string, 0, 4)
	if e.Title == "" {
		errMsgs = append(errMsgs, "title is empty")
	}
	if e.StartTime.Before(time.Now().Add(1 * time.Minute)) {
		errMsgs = append(errMsgs, "event time is expired")
	}
	if e.Duration <= 0 {
		errMsgs = append(errMsgs, "duration must be positive")
	}
	if e.UserID == "" {
		errMsgs = append(errMsgs, "user id is empty")
	}
	return errMsgs
}
