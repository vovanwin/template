package utils

import (
	"log"
	"math/rand"
	"time"
)

func ParseISO8601(iso8601 string) time.Time {
	if len(iso8601) < 3 {
		return time.Time{}
	}
	timeWant, err := time.Parse("2006-01-02T15:04:05Z", iso8601)
	if err != nil {
		log.Panic(err)
	}
	return timeWant
}

// "2006-01-02T15:04:05Z07:00"
func Parse3339(rfc3339 string) time.Time {
	if len(rfc3339) < 3 {
		return time.Time{}
	}
	timeWant, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		log.Panic(err)
	}
	return timeWant
}

// "2006-01-02T15:04:05Z07:00"
func UnixToParseISO8601(timestamp int) string {
	unixTimeUTC := time.Unix(1405544146, 0) //gives unix time stamp in utc

	unitTime := unixTimeUTC.Format("2006-01-02T15:04:05Z") // converts utc time to RFC3339 format
	return unitTime
}

func Randate() time.Time {
	min := time.Date(2023, 11, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2024, 4, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
