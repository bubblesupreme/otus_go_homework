package main

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	total   int64
	current int64
	str     string
}

func NewProgressBar(str string) *ProgressBar {
	return &ProgressBar{
		total:   0,
		current: 0,
		str:     str,
	}
}

func (pb *ProgressBar) Start(total int64) {
	pb.total = total
	pb.current = 0
	pb.redraw()
}

func (pb *ProgressBar) Finish() {
	pb.current = pb.total
	pb.redraw()
	fmt.Println("\nFinished")
}

func (pb *ProgressBar) Update(current int64) {
	pb.current = current
	pb.redraw()
}

func (pb *ProgressBar) UpdateStr(str string) {
	pb.str = str
	pb.redraw()
}

func (pb *ProgressBar) redraw() {
	percent := float64(pb.current) / float64(pb.total) * 100
	count := percent * 0.5
	fmt.Printf("\r[%-50s]%3.2f%% %8d/%d %s", strings.Repeat("#", int(count)), percent, pb.current, pb.total, pb.str)
}
