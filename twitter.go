package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

type Twitter struct {
}

func (Twitter) LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile(os.Getenv("HOME") + "/.turret/CREDENTIALS")
	if err != nil {
		return
	}
	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}
	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user)
	return
}

func (Twitter) Post(client *twittergo.Client, status string) {
	var (
		err   error
		req   *http.Request
		resp  *twittergo.APIResponse
		tweet *twittergo.Tweet
	)

	data := url.Values{}
	data.Set("status", status)
	body := strings.NewReader(data.Encode())
	req, err = http.NewRequest("POST", "/1.1/statuses/update.json", body)
	if err != nil {
		log.Printf("Error: Could not parse request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = client.SendRequest(req)
	if err != nil {
		log.Printf("Error: Could not send request: %v\n", err)
		return
	}
	tweet = &twittergo.Tweet{}
	err = resp.Parse(tweet)
	if err != nil {
		if rle, ok := err.(twittergo.RateLimitError); ok {
			log.Printf("Error: Rate limited, reset at %v\n", rle.Reset)
		} else if errs, ok := err.(twittergo.Errors); ok {
			for i, val := range errs.Errors() {
				log.Printf("Error #%v - ", i+1)
				log.Printf("Code: %v ", val.Code())
				log.Printf("Msg: %v\n", val.Message())
			}
		} else {
			log.Printf("Error: Problem parsing response: %v\n", err)
		}
		return
	}
	log.Printf("ID:                   %v\n", tweet.Id())
	log.Printf("Tweet:                %v\n", tweet.Text())
	log.Printf("User:                 %v\n", tweet.User().Name())

	if resp.HasRateLimit() {
		log.Printf("Rate limit:           %v\n", resp.RateLimit())
		log.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		log.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	}
}
