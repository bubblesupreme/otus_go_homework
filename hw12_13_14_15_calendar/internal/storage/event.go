package storage

import "time"

type Event struct {
	ID       int
	Title    string
	ClientID int
	Time     time.Time
}
