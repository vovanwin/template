package utils

import (
	"log"
	"time"
)

func ParseISO8601(iso8601 string) time.Time {
	timeWant, err := time.Parse("2006-01-02T15:04:05Z", iso8601)
	if err != nil {
		log.Panic(err)
	}
	return timeWant
}

// "2006-01-02T15:04:05Z07:00"
func Parse3339(rfc3339 string) time.Time {
	timeWant, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		log.Panic(err)
	}
	return timeWant
}
