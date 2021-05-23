package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrNotFoundEvent = errors.New("event not found")
	generator        = idGenerator{-1}
)

type memoryStorage struct {
	sync.RWMutex
	events []storage.Event
}

type checkTime struct {
	year  bool
	month bool
	day   bool
	week  bool
}

type idGenerator struct {
	counter int
}

func NewStorage() storage.Storage {
	return &memoryStorage{
		events: make([]storage.Event, 0),
	}
}

func (g *idGenerator) generateID() int {
	g.counter++
	return g.counter
}

func (s *memoryStorage) CreateEvent(ctx context.Context, e storage.Event) (int, error) {
	s.Lock()
	defer s.Unlock()

	select {
	case <-ctx.Done():
		return -1, storage.ErrContextDone
	default:
		e.ID = generator.generateID()
		s.events = append(s.events, e)

		return e.ID, nil
	}
}

func (s *memoryStorage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.Lock()
	defer s.Unlock()

	for i, e := range s.events {
		select {
		case <-ctx.Done():
			return storage.ErrContextDone
		default:
			if e.ID == event.ID {
				s.events[i] = event
				return nil
			}
		}
	}
	return ErrNotFoundEvent
}

func (s *memoryStorage) RemoveEvent(ctx context.Context, id int) error {
	s.Lock()
	defer s.Unlock()

	for i, e := range s.events {
		select {
		case <-ctx.Done():
			return storage.ErrContextDone
		default:
			if e.ID == id {
				s.events[i] = s.events[len(s.events)-1]
				s.events = s.events[:len(s.events)-1]
				return nil
			}
		}
	}
	return ErrNotFoundEvent
}

func (s *memoryStorage) GetDayEvents(ctx context.Context, eTime time.Time) ([]storage.Event, error) {
	return s.getTimeEvents(ctx, eTime, checkTime{day: true, week: true, month: true, year: true})
}

func (s *memoryStorage) GetWeekEvents(ctx context.Context, eTime time.Time) ([]storage.Event, error) {
	return s.getTimeEvents(ctx, eTime, checkTime{day: false, week: true, month: true, year: true})
}

func (s *memoryStorage) GetMonthEvents(ctx context.Context, eTime time.Time) ([]storage.Event, error) {
	return s.getTimeEvents(ctx, eTime, checkTime{day: false, week: false, month: true, year: true})
}

func (s *memoryStorage) getTimeEvents(ctx context.Context, eTime time.Time, t checkTime) ([]storage.Event, error) {
	s.Lock()
	defer s.Unlock()

	res := make([]storage.Event, 0)
	date := storage.ParseTime(eTime)
	for _, e := range s.events {
		select {
		case <-ctx.Done():
			return res, storage.ErrContextDone
		default:
			checkDate := storage.ParseTime(e.Time)
			if (t.year && checkDate.Year != date.Year) ||
				(t.month && checkDate.Month != date.Month) ||
				(t.week && checkDate.Week != date.Week) ||
				(t.day && checkDate.Day != date.Day) {
				continue
			}

			res = append(res, e)
		}
	}

	return res, nil
}
