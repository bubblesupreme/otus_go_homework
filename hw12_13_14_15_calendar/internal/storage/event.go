package storage

import "time"

type Event struct {
	ID       int32
	ClientID int32
	Title    string
	Time     time.Time
}
