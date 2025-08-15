package storage

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type StorageError struct {
	StatusCode      int
	ErrorMessage    string
	ConflictEventId string
	ConflictUserId  string
	Message         string
	Cause           error
	Zxx             *pgconn.PgError
}

func (se *StorageError) Error() string {
	return se.Message
}

func NewStorageError(message string) StorageError {
	return StorageError{
		Message: message,
	}
}

var (
	ErrEventIsNil          = errors.New("event is nil")
	ErrDateBusy            = errors.New("time is already taken by another event")
	ErrUserConflict        = errors.New("user is not the owner of the event")
	ErrEventNotFound       = errors.New("no events found")
	ErrEventAllreadyExists = errors.New("event with this id already exists")
)

const (
	ErrDateBusyMsgTemplate   = "time is already taken by another event: %s"
	ErrEventIdMsgTemplate    = "validate event id: %s"
	ErrEventIdWrapTemplate   = "validate event id: %w"
	ErrEmptyTitleMsg         = "title is empty"
	ErrEventTimeIsExpiredMsg = "event time is expired"
	ErrEmptyUserIdMsg        = "user id is empty"
)
