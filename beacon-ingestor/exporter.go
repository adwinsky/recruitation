package main

import (
	"log"
	"time"
)

func exportToGCP(ev Event) {
	time.Sleep(50 * time.Millisecond)

	log.Printf(
		"exported channel=%s type=%s ad=%s",
		ev.ChannelID,
		ev.EventType,
		ev.AdID,
	)
}
