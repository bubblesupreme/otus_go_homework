package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrStorageType = errors.New("storage type is not supported")
	ErrContextDone = errors.New("context is done")
)

type Storage interface {
	CreateEvent(ctx context.Context, e Event) (int, error) // returns id and error
	UpdateEvent(ctx context.Context, e Event) error
	RemoveEvent(ctx context.Context, id int) error
	GetDayEvents(ctx context.Context, eTime time.Time) ([]Event, error)
	GetWeekEvents(ctx context.Context, eTime time.Time) ([]Event, error)
	GetMonthEvents(ctx context.Context, eTime time.Time) ([]Event, error)
}
