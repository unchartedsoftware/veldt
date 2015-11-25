package twitter

import (
	"time"
)

func tweetDateToISO(tweetDate string) (string, error) {
	// attempt using this layout
	const layoutA = "Mon Jan 2 15:04:05 -0700 2006"
	tA, errA := time.Parse(layoutA, tweetDate)
	if errA == nil {
		return tA.Format(time.RFC3339), nil
	}
	// if it fails, attempt this layout
	const layoutB = "Mon Jan 2 15:04:05 MST 2006"
	tB, errB := time.Parse(layoutB, tweetDate)
	if errB != nil {
		return "", errB
	}
	return tB.Format(time.RFC3339), nil
}

func columnExists(col string) bool {
	if col != "" && col != "None" {
		return true
	}
	return false
}
