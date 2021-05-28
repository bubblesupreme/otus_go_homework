package memorystorage

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorageFull(t *testing.T) {
	s := NewStorage()

	ctx := context.Background()
	expectedDayIDs := make([]int32, 0)
	expectedWeekIDs := make([]int32, 0)
	expectedMonthIDs := make([]int32, 0)
	expectedMonthIDs2 := make([]int32, 0)

	id, err := s.CreateEvent(ctx, storage.Event{
		Title:    "Test 0 0",
		ClientID: 0,
		Time:     time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	})
	require.Nil(t, err)
	expectedDayIDs = append(expectedDayIDs, id)
	expectedWeekIDs = append(expectedWeekIDs, id)
	expectedMonthIDs = append(expectedMonthIDs, id)
	expectedMonthIDs2 = append(expectedMonthIDs2, id)

	id, err = s.CreateEvent(ctx, storage.Event{
		Title:    "Test 0 1",
		ClientID: 0,
		Time:     time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	})
	require.Nil(t, err)
	expectedDayIDs = append(expectedDayIDs, id)
	expectedWeekIDs = append(expectedWeekIDs, id)

	id, err = s.CreateEvent(ctx, storage.Event{
		Title:    "Test 0 2",
		ClientID: 0,
		Time:     time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC),
	})
	require.Nil(t, err)
	expectedWeekIDs = append(expectedWeekIDs, id)
	expectedMonthIDs = append(expectedMonthIDs, id)

	_, err = s.CreateEvent(ctx, storage.Event{
		Title:    "Test 0 3",
		ClientID: 0,
		Time:     time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
	})
	require.Nil(t, err)

	events, err := s.GetDayEvents(ctx, time.Date(2021, 2, 1, 5, 0, 0, 0, time.UTC))
	require.Nil(t, err)
	eventIds := make([]int32, len(events))
	for i, e := range events {
		eventIds[i] = e.ID
	}
	require.True(t, reflect.DeepEqual(eventIds, expectedDayIDs))

	events, err = s.GetWeekEvents(ctx, time.Date(2021, 2, 1, 5, 0, 0, 0, time.UTC))
	require.Nil(t, err)
	eventIds = make([]int32, len(events))
	for i, e := range events {
		eventIds[i] = e.ID
	}
	require.True(t, reflect.DeepEqual(eventIds, expectedWeekIDs))

	require.Nil(t, s.RemoveEvent(ctx, 1))

	events, err = s.GetMonthEvents(ctx, time.Date(2021, 2, 1, 5, 0, 0, 0, time.UTC))
	require.Nil(t, err)
	eventIds = make([]int32, len(events))
	for i, e := range events {
		eventIds[i] = e.ID
	}
	require.True(t, reflect.DeepEqual(eventIds, expectedMonthIDs))

	require.Nil(t, s.UpdateEvent(ctx, storage.Event{
		ID:       2,
		Title:    "Test 0 2",
		ClientID: 0,
		Time:     time.Date(2021, 3, 2, 0, 0, 0, 0, time.UTC),
	}))

	events, err = s.GetMonthEvents(ctx, time.Date(2021, 2, 1, 5, 0, 0, 0, time.UTC))
	require.Nil(t, err)
	eventIds = make([]int32, len(events))
	for i, e := range events {
		eventIds[i] = e.ID
	}
	require.True(t, reflect.DeepEqual(eventIds, expectedMonthIDs2))
}

func TestStorageNoExistRemove(t *testing.T) {
	s := NewStorage()

	ctx := context.Background()
	require.EqualError(t, s.RemoveEvent(ctx, 1), storage.ErrNotFoundEvent.Error())
}

func TestStorageNoExistUpdate(t *testing.T) {
	s := NewStorage()

	ctx := context.Background()
	require.EqualError(t, s.UpdateEvent(ctx, storage.Event{
		Title:    "Test 0 0",
		ClientID: 0,
		Time:     time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	}), storage.ErrNotFoundEvent.Error())
}

func TestStorageCancelContext(t *testing.T) {
	s := NewStorage()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := s.CreateEvent(ctx, storage.Event{
		Title:    "Test 0 0",
		ClientID: 0,
		Time:     time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	})
	require.EqualError(t, err, storage.ErrContextDone.Error())
}

func TestStorageCrowdedHours(t *testing.T) {
	s := NewStorage()

	ctx := context.Background()

	expectedIds := make([]int32, 1)
	var err error

	expectedIds[0], err = s.CreateEvent(ctx, storage.Event{
		Title:    "Test 0 0",
		ClientID: 0,
		Time:     time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	})
	require.Nil(t, err)

	_, err = s.CreateEvent(ctx, storage.Event{
		Title:    "Test 0 1",
		ClientID: 0,
		Time:     time.Date(2021, 2, 1, 50, 0, 0, 0, time.UTC),
	})
	require.Nil(t, err)

	events, err := s.GetDayEvents(ctx, time.Date(2021, 2, 1, 5, 0, 0, 0, time.UTC))
	require.Nil(t, err)
	eventIds := make([]int32, len(events))
	for i, e := range events {
		eventIds[i] = e.ID
	}
	require.True(t, reflect.DeepEqual(eventIds, expectedIds))
}
