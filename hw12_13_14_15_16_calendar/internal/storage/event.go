package storage

import (
	"fmt"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
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

// проверка двух эвентов на пересечение времени
func (e *Event) Overlaps(other *Event) bool {
	end1 := e.StartTime.Add(e.Duration)
	end2 := other.StartTime.Add(other.Duration)
	return e.StartTime.Before(end2) && end1.After(other.StartTime)
}

func (e *Event) GenerateEventId() {
	e.ID = uuid.New().String()
}

func ValidateEventId(eventId string) error {
	err := uuid.Validate(eventId)
	if err != nil {
		return NewSErrorWithCause(ErrFailedValidateEventIdTemplate, err)
	}
	return nil
}

func FullValidateEvent(e *Event) error {
	if e == nil {
		return ErrEventIsNil
	}
	errMsgs := make([]string, 0, 4)
	err := uuid.Validate(e.ID)
	if err != nil {
		errMsgs = append(errMsgs, fmt.Sprintf(ErrFailedValidateEventIdTemplate, err))
	}
	errMsgs = append(errMsgs, e.simpleValidate()...)
	return NewSErrorWithMsgArr(errMsgs...)
}

// без валидации ID
func ValidateEvent(e *Event) error {
	if e == nil {
		return ErrEventIsNil
	}
	errMsgs := e.simpleValidate()
	return NewSErrorWithMsgArr(errMsgs...)
}

func (e *Event) simpleValidate() []string {
	errMsgs := make([]string, 0, 3)
	if e.Title == "" {
		errMsgs = append(errMsgs, ErrEmptyTitleMsg)
	}
	if e.StartTime.Before(time.Now().Add(1 * time.Minute)) {
		errMsgs = append(errMsgs, ErrEventTimeIsExpiredMsg)
	}
	if e.UserID == "" {
		errMsgs = append(errMsgs, ErrEmptyUserIdMsg)
	}
	return errMsgs
}
