package storage

import (
	"fmt"
	"strings"
)

var (
	ErrEventIsNil         = NewSimpleSError("event is nil")
	ErrEventNotFound      = NewSimpleSError("event not found")
	ErrEventIDListIsEmpty = NewSimpleSError("event id list is empty")
)

const (
	ErrDateBusyMsgTemplate     = "time is already taken by another event: %s"
	ErrUserConflictMsgTemplate = "'%s' user is not the owner of the event, conflict with '%s'"
	ErrContextDone             = "context done"
)

type SError struct {
	Message string
	Cause   error
}

func (serr *SError) Error() string {
	if serr.Cause != nil {
		return fmt.Sprintf("%s: %v", serr.Message, serr.Cause)
	}
	return serr.Message
}

func (serr *SError) Unwrap() error {
	return serr.Cause
}

func NewSimpleSError(message string) error {
	return &SError{
		Message: message,
	}
}

func NewSError(message string, err error) error {
	return &SError{
		Message: message,
		Cause:   err,
	}
}

func NewSErrorWithTemplate(template string, messages ...any) error {
	if template == "" {
		return nil
	}
	return &SError{
		Message: fmt.Sprintf(template, messages...),
	}
}

func NewSErrorWithMsgArr(messages []string) error {
	message := joinString(messages)
	if message == "" {
		return nil
	}
	return &SError{
		Message: message,
	}
}

func NewSErrorWithCause(template string, err error) error {
	return &SError{
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
