package storage

import (
	"fmt"
	"strings"
)

var (
	ErrEventIsNil    = NewSimpleSError("event is nil")
	ErrEventNotFound = NewSimpleSError("event not found")
)

const (
	ErrDateBusyMsgTemplate     = "time is already taken by another event: %s"
	ErrUserConflictMsgTemplate = "user '%s' is not the owner of the event, conflict with '%s'"
	ErrContextDone             = "context done"
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

func (serr *StorageError) Unwrap() error {
	return serr.Cause
}

func NewSimpleSError(message string) error {
	return &StorageError{
		Message: message,
	}
}

func NewSError(message string, err error) error {
	return &StorageError{
		Message: message,
		Cause:   err,
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
