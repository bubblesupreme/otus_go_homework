package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const host = "0.beevik-ntp.pool.ntp.org"

func main() {
	fmt.Println("current time:", time.Now().Round(0))
	if t, err := ntp.Time(host); err != nil {
		log.Fatalf("failed to get time from %s: %s", host, err)
	} else {
		fmt.Println("exact time:", t.Round(0))
	}
}
