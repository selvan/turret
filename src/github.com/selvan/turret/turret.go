package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/selvan/turret/twitter"
)

var twitterClient, credentialErr = twitter.LoadCredentials()

var dateParseFormat = "2006-Jan-02"
var loc, _ = time.LoadLocation("Local")
var scheduledTweets []ScheduledTweet

type ScheduledTweet struct {
	Tweet       string
	ScheduledAt time.Time
}

type ByScheduledAt []ScheduledTweet

func (a ByScheduledAt) Len() int           { return len(a) }
func (a ByScheduledAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScheduledAt) Less(i, j int) bool { return a[i].ScheduledAt.Before(a[j].ScheduledAt) }

func ScheduledTweets() []ScheduledTweet {
	return scheduledTweets
}

func FilterByThreshold(thresholdTime time.Time, scheduledTweets []ScheduledTweet) (filteredTweets []ScheduledTweet, remainingTweets []ScheduledTweet) {
	for _i, _scheduledTweet := range scheduledTweets {
		if thresholdTime.After(_scheduledTweet.ScheduledAt) {
			continue
		}
		filteredTweets = scheduledTweets[:_i]
		remainingTweets = scheduledTweets[_i:]
		break
	}
	return
}

func ParseCSV(cutoffTime time.Time, reader *csv.Reader) {
	scheduledTweets = []ScheduledTweet{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		dateTime := strings.Split(record[0], " ")
		dateParts, _ := time.Parse(dateParseFormat, dateTime[0])
		timeParts := strings.Split(dateTime[1], ":")

		hour, _ := strconv.Atoi(timeParts[0])
		minute, _ := strconv.Atoi(timeParts[1])

		schduledAt := time.Date(dateParts.Year(), dateParts.Month(), dateParts.Day(), hour, minute, 0, 0, loc)

		// Ignore any old tweets
		if schduledAt.Before(cutoffTime) {
			log.Printf("Ignoring : Today's date (%s) > schedule date (%s) : %s", cutoffTime.Format(dateParseFormat), schduledAt.Format(dateParseFormat), record)
			continue
		}

		tweet := record[1]
		scheduledTweets = append(scheduledTweets, ScheduledTweet{tweet, schduledAt})
	}
	sort.Sort(ByScheduledAt(scheduledTweets))
	return
}

func postTweets(tweets []ScheduledTweet) {
	if tweets != nil && len(tweets) > 0 {
		log.Printf("Number of tweets ready to post is %d ", len(tweets))
	}

	for i, tweet := range tweets {
		if i > 0 {
			// Delay by 10 seconds, if there is more than 1 tweet
			time.Sleep(10 * time.Second)
		}
		twitter.Post(twitterClient, tweet.Tweet)
	}
}

func startClock() chan bool {
	stopHandle := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				filteredTweets, remainingTweets := FilterByThreshold(time.Now(), scheduledTweets)
				scheduledTweets = remainingTweets
				if (len(filteredTweets) == 0) && (len(scheduledTweets) == 0) {
					log.Fatal("No more scheduled tweets and no more tweets to post, exiting.")
				}
				go postTweets(filteredTweets)
			case <-stopHandle:
				return
			}
		}
	}()

	return stopHandle
}

func main() {

	if credentialErr != nil {
		log.Fatalf("Could not parse  ~/.turret/CREDENTIALS file: %v\n", credentialErr)
	}

	csvFile, csvFileError := os.Open("tweets.csv")
	if csvFileError != nil {
		log.Fatalf("Could not parse  tweets.csv file: %v\n", csvFileError)
	}

	ParseCSV(time.Now(), csv.NewReader(csvFile))
	<-startClock()
}
