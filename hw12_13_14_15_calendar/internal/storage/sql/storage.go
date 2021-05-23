package sqlstorage

import (
	"context"
	"fmt"
	"time"

	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v4"
)

const (
	driver = "postgres"
	table  = "events"
)

type sqlStorage struct {
	login    string
	password string
	port     int
	host     string
	dbname   string
	db       *pgx.Conn
}

func NewStorage(login string, password string, port int, host string, dbname string) (storage.Storage, error) {
	storage := &sqlStorage{
		login:    login,
		password: password,
		port:     port,
		dbname:   dbname,
		host:     host,
	}

	ctx := context.Background()
	if err := storage.Connect(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if err := storage.Close(ctx); err != nil {
			log.Fatal(fmt.Sprintf("failed to close connection: %s", err), log.Fields{})
		}
	}()

	return storage, nil
}

func (s *sqlStorage) Connect(ctx context.Context) error {
	var err error
	s.db, err = pgx.Connect(ctx, fmt.Sprintf("%s://%s:%s@%s:%d/%s", driver, s.login, s.password, s.host, s.port, s.dbname))

	return err
}

func (s *sqlStorage) Close(ctx context.Context) error {
	return s.db.Close(ctx)
}

func (s *sqlStorage) CreateEvent(ctx context.Context, e storage.Event) (int, error) {
	if err := s.Connect(ctx); err != nil {
		return -1, err
	}
	defer func() {
		if err := s.Close(ctx); err != nil {
			log.Fatal(fmt.Sprintf("failed to close connection: %s", err), log.Fields{})
		}
	}()

	d := storage.ParseTime(e.Time)
	id := -1
	err := s.db.QueryRow(ctx,
		`
INSERT 
INTO $1 
(title, client_id, "year", "month", "day", "hour", "minutes", "seconds") 
VALUES ("$2", $3, $4, $5, $6, $7, $8, $9) 
RETURNING id;`, table, e.Title, e.ClientID, d.Year, d.Month, d.Day, d.Hour, d.Minutes, d.Seconds).Scan(&id)

	return id, err
}

func (s *sqlStorage) UpdateEvent(ctx context.Context, e storage.Event) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}
	defer func() {
		if err := s.Close(ctx); err != nil {
			log.Fatal(fmt.Sprintf("failed to close connection: %s", err), log.Fields{})
		}
	}()

	date := storage.ParseTime(e.Time)
	_, err := s.db.Query(ctx,
		`
UPDATE $1 
SET title = $2,
    client_id = $3,
    "year" = $4,
    "month" = $5,
    "day" = $6,
    "hour" = $7,
    "minutes" = $8,
    "seconds" = $9
WHERE id = $10;`,
		table, e.Title, e.ClientID, date.Year, int(date.Month), date.Day, date.Hour, date.Minutes, date.Seconds, e.ID)

	return err
}

func (s *sqlStorage) RemoveEvent(ctx context.Context, id int) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}
	defer func() {
		if err := s.Close(ctx); err != nil {
			log.Fatal(fmt.Sprintf("failed to close connection: %s", err), log.Fields{})
		}
	}()

	_, err := s.db.Query(ctx,
		`
DELETE 
FROM $1 
WHERE id = $2;`, table, id)
	return err
}

func (s *sqlStorage) GetDayEvents(ctx context.Context, eTime time.Time) ([]storage.Event, error) {
	date := storage.ParseTime(eTime)
	rows, err := s.db.Query(ctx,
		`
SELECT * 
FROM $1 
WHERE ("year" = $2 AND "month" = $3 AND "day" = $4);`, table, date.Year, int(date.Month), date.Day)
	if err != nil {
		return nil, err
	}

	return rowsToEventSlice(rows)
}

func (s *sqlStorage) GetWeekEvents(ctx context.Context, eTime time.Time) ([]storage.Event, error) {
	if err := s.Connect(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if err := s.Close(ctx); err != nil {
			log.Fatal(fmt.Sprintf("failed to close connection: %s", err), log.Fields{})
		}
	}()

	date := storage.ParseTime(eTime)
	startDay := date.Day - date.Week + 1
	endDay := date.Day + (7 - date.Week)

	rows, err := s.db.Query(ctx,
		`
SELECT * 
FROM $1 
WHERE ("year" = $2 AND "month" = $3 AND "day" >= $4 AND "day" <= $5);`, table, date.Year, int(date.Month), startDay, endDay)
	if err != nil {
		return nil, err
	}

	return rowsToEventSlice(rows)
}

func (s *sqlStorage) GetMonthEvents(ctx context.Context, eTime time.Time) ([]storage.Event, error) {
	if err := s.Connect(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if err := s.Close(ctx); err != nil {
			log.Fatal(fmt.Sprintf("failed to close connection: %s", err), log.Fields{})
		}
	}()

	date := storage.ParseTime(eTime)
	rows, err := s.db.Query(ctx,
		`
SELECT * 
FROM $1 
WHERE ("year" = $2 AND "month" = $3);`, table, date.Year, int(date.Month))
	if err != nil {
		return nil, err
	}

	return rowsToEventSlice(rows)
}

func rowsToEventSlice(rows pgx.Rows) ([]storage.Event, error) {
	res := make([]storage.Event, 0)
	for rows.Next() {
		e := storage.Event{}
		d := storage.Date{}
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.ClientID,
			&d.Year,
			&d.Month,
			&d.Day,
			&d.Hour,
			&d.Minutes,
			&d.Seconds,
		); err != nil {
			return res, err
		}

		e.Time = storage.DateToTime(d)
		res = append(res, e)
	}

	return res, rows.Err()
}
