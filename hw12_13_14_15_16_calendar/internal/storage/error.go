package storage

import "errors"

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
