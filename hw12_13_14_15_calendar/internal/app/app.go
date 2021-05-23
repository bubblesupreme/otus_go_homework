package app

import (
	"context"
	"io"

	"google.golang.org/protobuf/types/known/timestamppb"

	eventspb "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/api"

	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	eventspb.UnimplementedEventServiceServer
	storage storage.Storage
}

type Logger interface {
	SetLevel(levelStr string)
	SetOutput(output io.Writer)
	Info(msg string, f log.Fields)
	Error(msg string, f log.Fields)
	Warning(msg string, f log.Fields)
	Fatal(msg string, f log.Fields)
}

func NewApp(storage storage.Storage) *App {
	return &App{storage: storage}
}

func (a App) CreateEvent(ctx context.Context, e *eventspb.Event) (*eventspb.EventID, error) {
	res := eventspb.EventID{}
	var err error
	res.ID, err = a.storage.CreateEvent(ctx, storage.Event{
		Title:    e.GetTitle(),
		ClientID: e.GetClientID(),
		Time:     e.GetTime().AsTime(),
	})
	return &res, err
}

func (a App) UpdateEvent(ctx context.Context, e *eventspb.Event) (*eventspb.Empty, error) {
	return &eventspb.Empty{}, a.storage.UpdateEvent(ctx, storage.Event{
		ID:       e.GetID(),
		Title:    e.GetTitle(),
		ClientID: e.GetClientID(),
		Time:     e.GetTime().AsTime(),
	})
}

func (a App) RemoveEvent(ctx context.Context, e *eventspb.EventID) (*eventspb.Empty, error) {
	return &eventspb.Empty{}, a.storage.RemoveEvent(ctx, e.ID)
}

func (a App) GetDayEvents(ctx context.Context, e *eventspb.Time) (*eventspb.Events, error) {
	events, err := a.storage.GetDayEvents(ctx, e.GetTime().AsTime())
	return eventsStorageToPB(events), err
}

func (a App) GetWeekEvents(ctx context.Context, e *eventspb.Time) (*eventspb.Events, error) {
	events, err := a.storage.GetWeekEvents(ctx, e.GetTime().AsTime())
	return eventsStorageToPB(events), err
}

func (a App) GetMonthEvents(ctx context.Context, e *eventspb.Time) (*eventspb.Events, error) {
	events, err := a.storage.GetMonthEvents(ctx, e.GetTime().AsTime())
	return eventsStorageToPB(events), err
}

func eventsStorageToPB(events []storage.Event) *eventspb.Events {
	res := eventspb.Events{}
	res.Events = make([]*eventspb.Event, len(events))
	for i, e := range events {
		res.Events[i] = &eventspb.Event{
			Time:     timestamppb.New(e.Time),
			ID:       e.ID,
			Title:    e.Title,
			ClientID: e.ClientID,
		}
	}
	return &res
}
