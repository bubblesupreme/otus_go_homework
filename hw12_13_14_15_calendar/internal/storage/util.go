package storage

import "time"

type Date struct {
	Year    int
	Month   time.Month
	Week    int
	Day     int
	Hour    int
	Minutes int
	Seconds int
}

func ParseTime(t time.Time) Date {
	var res Date
	res.Year, res.Month, res.Day = t.Date()
	res.Hour, res.Minutes, res.Seconds = t.Clock()
	_, res.Week = t.ISOWeek()

	return res
}

func DateToTime(d Date) time.Time {
	return time.Date(d.Year, d.Month, d.Day, d.Hour, d.Minutes, d.Seconds, 0, time.UTC)
}
