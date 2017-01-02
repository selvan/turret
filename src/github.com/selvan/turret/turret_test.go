package main

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"
)

func testEq(a, b []ScheduledTweet) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestFilterByThreshold(t *testing.T) {
	loc, _ := time.LoadLocation("Local")
	scheduledTweets := []ScheduledTweet{
		ScheduledTweet{Tweet: "One", ScheduledAt: time.Date(2017, time.Month(01), 01, 9, 1, 0, 0, loc)},
		ScheduledTweet{Tweet: "Two", ScheduledAt: time.Date(2017, time.Month(01), 01, 9, 5, 0, 0, loc)},
		ScheduledTweet{Tweet: "Three", ScheduledAt: time.Date(2017, time.Month(01), 01, 9, 7, 0, 0, loc)},
	}
	filteredTweets, remainingTweets := FilterByThreshold(time.Date(2017, time.Month(01), 01, 9, 6, 0, 0, loc), scheduledTweets)
	if !testEq(filteredTweets, scheduledTweets[:2]) {
		t.Error("Filtered tweets are not good")
	}

	if !testEq(remainingTweets, scheduledTweets[2:]) {
		t.Error("Remaining tweets are not good")
	}
}

func TestParseCSV(t *testing.T) {
	csvData := `
2117-Dec-01 10:30,"Hello y"
2117-Dec-01 11:30,"Hello z"
2117-Dec-01 8:30,"Hello x"
`
	reader := csv.NewReader(strings.NewReader(csvData))
	now := time.Now()
	ParseCSV(now, reader)
	scheduledTweets := ScheduledTweets()

	if !(len(scheduledTweets) == 3) {
		t.Error("Wrong number of tweets?")
	}

	if !(scheduledTweets[0].Tweet == "Hello x") {
		t.Error("Not sorted by time ? ", scheduledTweets[0])
	}

	if !(scheduledTweets[1].Tweet == "Hello y") {
		t.Error("Not sorted by time ? ", scheduledTweets[1])
	}

	if !(scheduledTweets[2].Tweet == "Hello z") {
		t.Error("Not sorted by time ? ", scheduledTweets[2])
	}
}

func TestParseCSVIgnoreOldTweets(t *testing.T) {
	csvData := `
2016-Dec-01 10:30,"Hello y"
2117-Dec-01 11:30,"Hello z"
2016-Dec-01 8:30,"Hello x"
`
	reader := csv.NewReader(strings.NewReader(csvData))
	now := time.Now()
	ParseCSV(now, reader)
	scheduledTweets := ScheduledTweets()

	if !(len(scheduledTweets) == 1) {
		t.Error("Not filtered ?. Old tweets are still present.")
	}

	if !(scheduledTweets[0].Tweet == "Hello z") {
		t.Error("Not sorted by time ? ", scheduledTweets[0])
	}
}
