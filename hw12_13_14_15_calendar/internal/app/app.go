package app

import (
	"context"
	"io"
	"time"

	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
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

func New(storage storage.Storage) *App {
	return &App{storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, title string, clientID int, eTime time.Time) (id int, err error) {
	return a.storage.CreateEvent(ctx, storage.Event{
		Title:    title,
		ClientID: clientID,
		Time:     eTime,
	})
}
