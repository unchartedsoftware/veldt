package twitter

import (
	"fmt"
	"time"
)

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
