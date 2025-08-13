package storage

import (
	"errors"
	"strings"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

var (
	ErrDateBusy             = errors.New("time is already taken by another event")
	ErrNotifyTooLate        = errors.New("notification time has already expired")
	ErrUserConflict         = errors.New("user is not the owner of the event")
	ErrEventNotFound        = errors.New("no events found")
	ErrEventAllreadyCreated = errors.New("event with this ID has already been created")
	//ErrEmptyTitle           = errors.New("title is empty")
)

const (
	ErrEmptyTitleMsg         = "title is empty"
	ErrEventTimeIsExpiredMsg = "event time is expired"
	ErrEmptyDescriptionMsg   = "description is empty"
	ErrEmptyUserIdMsg        = "user id is empty"
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

type Storage interface {
	Add(eventRef *Event) error
	Update(id string, eventRef *Event) error
	Delete(id string) error
	ListDay(day time.Time) ([]Event, error)
	ListWeek(start time.Time) ([]Event, error)
	ListMonth(start time.Time) ([]Event, error)
	Get(id string) (Event, error)
	ListPeriodByUserId(start time.Time, duration time.Duration, userId string)
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

func (e *Event) Validate() error {
	errMsgs := make([]string, 0, 4)
	if e.Title == "" {
		errMsgs = append(errMsgs, ErrEmptyTitleMsg)
	}
	if e.StartTime.Before(time.Now()) {
		errMsgs = append(errMsgs, ErrEventTimeIsExpiredMsg)
	}
	if e.Description == "" {
		errMsgs = append(errMsgs, ErrEmptyDescriptionMsg)
	}
	if e.UserID == "" {
		errMsgs = append(errMsgs, ErrEmptyUserIdMsg)
	}

	errMsg := JoinWithComma(errMsgs)
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}

func (e *Event) CheckAndGenerateId() {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
}

func JoinWithComma(items []string) string {
	var nonEmpty []string
	for _, item := range items {
		if item != "" {
			nonEmpty = append(nonEmpty, item)
		}
	}
	return strings.Join(nonEmpty, ", ")
}
