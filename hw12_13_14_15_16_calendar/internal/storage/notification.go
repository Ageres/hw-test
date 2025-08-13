package storage

import "time"

type Notification struct {
	ID        string
	Title     string
	StartTime time.Time
	UserID    string
}
