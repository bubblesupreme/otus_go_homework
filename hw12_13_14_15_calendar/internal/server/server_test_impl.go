package server

import (
	"context"
	"strconv"
	"testing"
	"time"

	eventspb "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/api"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TestClient interface {
	CreateEvent(ctx context.Context, in *eventspb.Event) (*eventspb.EventID, error)
	UpdateEvent(ctx context.Context, in *eventspb.Event) (*eventspb.Empty, error)
	RemoveEvent(ctx context.Context, in *eventspb.EventID) (*eventspb.Empty, error)
	GetDayEvents(ctx context.Context, in *eventspb.Time) (*eventspb.Events, error)
	GetWeekEvents(ctx context.Context, in *eventspb.Time) (*eventspb.Events, error)
	GetMonthEvents(ctx context.Context, in *eventspb.Time) (*eventspb.Events, error)
}

func init() {
	viper.AutomaticEnv()
	viper.BindEnv("dblogin", "POSTGRES_USER")
	viper.BindEnv("dbname", "POSTGRES_DB")
	viper.BindEnv("dbpassword", "POSTGRES_PASSWORD")
	viper.BindEnv("dbport", "POSTGRES_PORT")
	viper.BindEnv("dbhost", "POSTGRES_HOST")
}

func GetMemoryStorage() (storage.Storage, error) {
	return memorystorage.NewStorage(), nil
}

func GetSQLStorage() (storage.Storage, error) {
	port, err := strconv.Atoi(viper.Get("dbport").(string))
	if err != nil {
		return nil, err
	}
	return sqlstorage.NewStorage(viper.Get("dblogin").(string), viper.Get("dbpassword").(string),
		port, viper.Get("dbhost").(string), viper.Get("dbname").(string))
}

func EmptyStorageImpl(t *testing.T, client TestClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	date := time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)
	timestmp := eventspb.Time{
		Time: timestamppb.New(date),
	}

	events, err := client.GetDayEvents(ctx, &timestmp)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events.Events))

	events, err = client.GetWeekEvents(ctx, &timestmp)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events.Events))

	events, err = client.GetMonthEvents(ctx, &timestmp)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events.Events))

	_, err = client.UpdateEvent(ctx, &eventspb.Event{
		ClientID: 0,
		ID:       0,
		Title:    "",
		Time:     timestmp.Time,
	})
	assert.Error(t, err)

	_, err = client.RemoveEvent(ctx, &eventspb.EventID{
		ID: 0,
	})
	assert.Error(t, err)
}

func CommonImpl(t *testing.T, client TestClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	date0 := time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)
	expEvents := make([]eventspb.Event, 2)
	expEvents[0] = eventspb.Event{
		ID:       -1,
		Title:    "first test event",
		ClientID: 2,
		Time:     timestamppb.New(date0),
	}
	expEvents[1] = eventspb.Event{
		ID:       -1,
		Title:    "second test event",
		ClientID: 2,
		Time:     timestamppb.New(time.Date(2021, 3, 3, 0, 0, 0, 0, time.UTC)),
	}

	id, err := client.CreateEvent(ctx, &expEvents[0])
	assert.NoError(t, err)
	assert.True(t, id.ID != expEvents[0].ID)
	expEvents[0].ID = id.ID

	id, err = client.CreateEvent(ctx, &expEvents[1])
	assert.NoError(t, err)
	assert.True(t, id.ID != expEvents[1].ID)
	expEvents[1].ID = id.ID

	timestmp := eventspb.Time{
		Time: timestamppb.New(date0),
	}
	events, err := client.GetDayEvents(ctx, &timestmp)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(events.Events))
	compareEvent(t, &expEvents[0], events.GetEvents()[0])

	events, err = client.GetWeekEvents(ctx, &timestmp)
	assert.NoError(t, err)
	assert.Equal(t, len(expEvents), len(events.Events))
	compareEvent(t, &expEvents[0], events.GetEvents()[0])
	compareEvent(t, &expEvents[1], events.GetEvents()[1])

	events, err = client.GetMonthEvents(ctx, &timestmp)
	assert.NoError(t, err)
	assert.Equal(t, len(expEvents), len(events.Events))
	compareEvent(t, &expEvents[0], events.GetEvents()[0])
	compareEvent(t, &expEvents[1], events.GetEvents()[1])

	expEvents[0].Time = timestamppb.New(time.Date(2020, 3, 3, 0, 0, 0, 0, time.UTC))
	_, err = client.UpdateEvent(ctx, &expEvents[0])
	assert.NoError(t, err)

	events, err = client.GetWeekEvents(ctx, &timestmp)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(events.Events))
	compareEvent(t, &expEvents[1], events.GetEvents()[0])

	_, err = client.RemoveEvent(ctx, &eventspb.EventID{
		ID: expEvents[0].ID,
	})
	assert.NoError(t, err)

	_, err = client.RemoveEvent(ctx, &eventspb.EventID{
		ID: expEvents[1].ID,
	})
	assert.NoError(t, err)
}

func compareEvent(t *testing.T, expected, actual *eventspb.Event) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Time.AsTime(), actual.Time.AsTime())
	assert.Equal(t, expected.ClientID, actual.ClientID)
	assert.Equal(t, expected.Title, actual.Title)
}
