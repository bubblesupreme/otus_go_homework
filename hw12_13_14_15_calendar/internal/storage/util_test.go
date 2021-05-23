package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	year    = 2021
	month   = 2
	week    = 6
	day     = 8
	hour    = 5
	minutes = 46
	seconds = 13
)

func TestUtilParseTime(t *testing.T) {
	checkTime := time.Date(year, month, day, hour, minutes, seconds, 0, time.UTC)
	d := ParseTime(checkTime)

	require.Equal(t, year, d.Year)
	require.Equal(t, time.Month(month), d.Month)
	require.Equal(t, week, d.Week)
	require.Equal(t, day, d.Day)
	require.Equal(t, hour, d.Hour)
	require.Equal(t, minutes, d.Minutes)
	require.Equal(t, seconds, d.Seconds)
}

func TestUtilToDate(t *testing.T) {
	d := Date{year, month, week, day, hour, minutes, seconds}
	checkTime := DateToTime(d)
	require.Equal(t, time.Date(year, month, day, hour, minutes, seconds, 0, time.UTC), checkTime)
}
