package storage

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEventIsNil    = NewSimpleSError("event is nil")
	ErrEventNotFound = NewSimpleSError("events not found")
	ErrDateBusy      = errors.New("time is already taken by another event")
	ErrUserConflict  = errors.New("user is not the owner of the event")

	ErrEventAllreadyExists = errors.New("event with this id already exists")
)

const (
	ErrDateBusyMsgTemplate           = "time is already taken by another event: %s"
	ErrFailedValidateEventIdTemplate = "failed to validate event id: %v"
	ErrDatabaseTimeoutMsgTemplate    = "database timeout: %s"
	ErrDatabaseMsgTemplate           = "database error: %s"
	ErrUserConflictMsgTemplate       = "user '%s' is not the owner of the event, conflict with '%s'"
	ErrFailedAddEventTemplate        = "failed to add event: %v"
	ErrFailedUpdateEventTemplate     = "failed to update event: %v"
	ErrFailedDeleteEventTemplate     = "failed to delete event: %v"
	ErrFailedListEventTemplate       = "failed to list event: %v"
)

const (
	ErrEmptyTitleMsg         = "title is empty"
	ErrEventTimeIsExpiredMsg = "event time is expired"
	ErrEmptyUserIdMsg        = "user id is empty"
	ErrEventNotFoundMsg      = "no events found"
)

type StorageError struct {
	StatusCode      int
	ErrorMessage    string
	ConflictEventId string
	ConflictUserId  string
	Message         string
	Cause           error
}

func (serr *StorageError) Error() string {
	return serr.Message
}

func NewSimpleSError(message string) error {
	return &StorageError{
		Message: message,
	}
}

func NewSErrorWithTemplate(template string, messages ...string) error {
	return &StorageError{
		Message: fmt.Sprintf(template, messages),
	}
}

func NewSErrorWithMsgArr(messages ...string) error {
	message := joinString(messages)
	if message == "" {
		return nil
	}
	return &StorageError{
		Message: message,
	}
}

func NewSErrorWithCause(template string, err error) error {
	return &StorageError{
		Message: fmt.Sprintf(template, err),
		Cause:   err,
	}
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
