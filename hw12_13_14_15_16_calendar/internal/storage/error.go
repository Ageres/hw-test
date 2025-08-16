package storage

import (
	"fmt"
	"strings"
)

var (
	ErrEventIsNil    = NewSError("event is nil")
	ErrEventNotFound = NewSError("event not found")
)

const (
	ErrDateBusyMsgTemplate           = "time is already taken by another event: %s"
	ErrFailedValidateEventIdTemplate = "failed to validate event id: %v"
	ErrUserConflictMsgTemplate       = "user '%s' is not the owner of the event, conflict with '%s'"
	ErrContextDoneTemplate           = "context done: '%s'"
)

type StorageError struct {
	Message string
	Cause   error
}

func (serr *StorageError) Error() string {
	if serr.Cause != nil {
		return fmt.Sprintf("%s: %v", serr.Message, serr.Cause)
	}
	return serr.Message
}

func NewSError(message string) error {
	return &StorageError{
		Message: message,
	}
}

func NewSErrorWithTemplate(template string, messages ...any) error {
	if template == "" {
		return nil
	}
	return &StorageError{
		Message: fmt.Sprintf(template, messages...),
	}
}

func NewSErrorWithMsgArr(messages []string) error {
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
	if len(items) == 0 {
		return ""
	}

	var nonEmpty []string
	for _, item := range items {
		if item != "" {
			nonEmpty = append(nonEmpty, item)
		}
	}

	if len(nonEmpty) == 0 {
		return ""
	}
	return strings.Join(nonEmpty, "; ")
}
