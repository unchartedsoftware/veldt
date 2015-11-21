package twitter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/unchartedsoftware/prism/binning"
)

// TweetData represents a single row of twitter TSV data.
type TweetData struct {
	ISODate   string
	UserID    string
	Username  string
	TweetID   string
	TweetText string
	Hashtags  []string
	LonLat    *binning.LonLat
	Country   *string
	Location  *string
	City      *string
	Language  *string
}

func tweetDateToISO(tweetDate string) string {
	const layout = "Mon Jan 2 15:04:05 -0700 2006"
	t, err := time.Parse(layout, tweetDate)
	if err != nil {
		fmt.Println("Error parsing date: " + tweetDate)
		return ""
	}
	return t.Format(time.RFC3339)
}

func columnExists(col string) bool {
	if col != "" && col != "None" {
		return true
	}
	return false
}

// ParseTweetData transforms a csv string into a tweet data struct.
func ParseTweetData(tweetCsv []string) *TweetData {
	// CSV line as array:
	//     0: 'Fri Jan 04 18:42:42 +0000 2013',
	//     1: '242573761',
	//     2: 'AdioAsh5',
	//     3:  '287267829735100416',
	//     4:  "Blah blah blah blah blah",
	//     5:  '',
	//     6:  '-73.94068643', {lon}
	//     7:  '40.66179087', {lat}
	//     8:  'United States',
	//     9:  'New York, NY',
	//     10:  'city',
	//     11:  'en'
	tweet := TweetData{
		ISODate:   tweetDateToISO(tweetCsv[0]),
		UserID:    tweetCsv[1],
		Username:  tweetCsv[2],
		TweetID:   tweetCsv[3],
		TweetText: tweetCsv[4],
		Hashtags:  make([]string, 0),
	}
	// lon / lat data may not exist
	if columnExists(tweetCsv[6]) && columnExists(tweetCsv[7]) {
		lon, lonErr := strconv.ParseFloat(tweetCsv[6], 64)
		lat, latErr := strconv.ParseFloat(tweetCsv[7], 64)
		if lonErr == nil && latErr == nil {
			tweet.LonLat = &binning.LonLat{
				Lat: lat,
				Lon: lon,
			}
		}
	}
	// hashtags may not exist
	if columnExists(tweetCsv[5]) {
		tweet.Hashtags = strings.Split(strings.TrimSpace(tweetCsv[5]), "#")
	}
	// country data may not exist
	if columnExists(tweetCsv[8]) {
		tweet.Country = &tweetCsv[8]
	}
	// location data may not exist
	if columnExists(tweetCsv[9]) {
		tweet.Location = &tweetCsv[9]
	}
	// city data may not exist
	if columnExists(tweetCsv[10]) {
		tweet.City = &tweetCsv[10]
	}
	// language data may not exist
	if columnExists(tweetCsv[11]) {
		tweet.Language = &tweetCsv[11]
	}
	return &tweet
}
