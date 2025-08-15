package storage

import (
	"errors"
	"fmt"
	"strings"
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

func (e *Event) Overlaps(other *Event) bool {
	end1 := e.StartTime.Add(e.Duration)
	end2 := other.StartTime.Add(other.Duration)
	return e.StartTime.Before(end2) && end1.After(other.StartTime)
}

func (e *Event) FullValidate() error {
	errMsgs := make([]string, 0, 4)

	err := uuid.Validate(e.ID)
	if err != nil {
		errMsgs = append(errMsgs, fmt.Sprintf(ErrEventIdMsgTemplate, err))
	}
	if e.Title == "" {
		errMsgs = append(errMsgs, ErrEmptyTitleMsg)
	}
	if e.StartTime.Before(time.Now()) {
		errMsgs = append(errMsgs, ErrEventTimeIsExpiredMsg)
	}
	if e.UserID == "" {
		errMsgs = append(errMsgs, ErrEmptyUserIdMsg)
	}

	errMsg := joinString(errMsgs)
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}

// без валидации ID
func (e *Event) Validate() error {
	errMsgs := make([]string, 0, 3)
	if e.Title == "" {
		errMsgs = append(errMsgs, ErrEmptyTitleMsg)
	}
	if e.StartTime.Before(time.Now()) {
		errMsgs = append(errMsgs, ErrEventTimeIsExpiredMsg)
	}
	if e.UserID == "" {
		errMsgs = append(errMsgs, ErrEmptyUserIdMsg)
	}

	errMsg := joinString(errMsgs)
	if errMsg != "" {
		return &StorageError{
			StatusCode: 400,
			Message:    errMsg,
		}
	}
	return nil
}

func (e *Event) GenerateId() {
	e.ID = uuid.New().String()
}

func ValidateEventNotNil(e *Event) error {
	if e == nil {
		return ErrEventIsNil
	}
	return nil
}

func ValidateEventId(eventId string) error {
	err := uuid.Validate(eventId)
	if err != nil {
		return fmt.Errorf(ErrEventIdWrapTemplate, err)
	}
	return nil
}

func joinString(items []string) string {
	var nonEmpty []string
	for _, item := range items {
		if item != "" {
			nonEmpty = append(nonEmpty, item)
		}
	}
	return strings.Join(nonEmpty, "; ")
}
