package storage

import "time"

type Notification struct {
	ID        string //UUID, ID из Event
	Title     string
	StartDate time.Time
	UserID    string //UUID, UserID из Event
}
