package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrContextDone         = errors.New("context is done")
	ErrNotFoundEvent       = errors.New("event not found")
	ErrFoundMultipleEvents = errors.New("multiple events were found")
)

type Storage interface {
	CreateEvent(ctx context.Context, e Event) (int32, error) // returns id and error
	UpdateEvent(ctx context.Context, e Event) error
	RemoveEvent(ctx context.Context, id int32) error
	GetDayEvents(ctx context.Context, eTime time.Time) ([]Event, error)
	GetWeekEvents(ctx context.Context, eTime time.Time) ([]Event, error)
	GetMonthEvents(ctx context.Context, eTime time.Time) ([]Event, error)
}
