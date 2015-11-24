package progress

import (
	"fmt"
	"time"

	"github.com/unchartedsoftware/prism/util"
)

// Time represents the hour, minutes, and seconds of a given time.
type Time struct {
	Seconds uint64
	Minutes uint64
	Hours   uint64
}

func formatTime(totalSeconds uint64) Time {
	totalMinutes := totalSeconds / 60
	seconds := totalSeconds % 60
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return Time{
		Seconds: seconds,
		Minutes: minutes,
		Hours:   hours,
	}
}

func getTimestamp() uint64 {
	return uint64(time.Now().Unix())
}

var startTime uint64

// PrintProgress will print a human readable progress message for a given task.
func PrintProgress(totalBytes uint64, bytes uint64) {
	if startTime == 0 {
		startTime = getTimestamp()
	}
	elapsed := getTimestamp() - startTime
	percentComplete := 100 * (float64(bytes) / float64(totalBytes))
	bytesPerSecond := 1.0
	if elapsed > 0 {
		bytesPerSecond = float64(bytes) / float64(elapsed)
	}
	estimatedSecondsRemaining := (float64(totalBytes) - float64(bytes)) / bytesPerSecond
	formattedTime := formatTime(uint64(estimatedSecondsRemaining))
	formattedBytes := util.FormatBytes(float64(bytes))
	formattedBytesPerSecond := util.FormatBytes(bytesPerSecond)
	fmt.Printf("\rProcessed %+9s at %+9sps, %6.2f%% complete, estimated time remaining: %2d:%02d:%02d",
		formattedBytes,
		formattedBytesPerSecond,
		percentComplete,
		formattedTime.Hours%100,
		formattedTime.Minutes,
		formattedTime.Seconds)
}

// PrintTotalDuration prints the total duration of the processed task.
func PrintTotalDuration() {
	formattedTime := formatTime(getTimestamp() - startTime)
	fmt.Printf("\nTask completed in %d:%02d:%02d\n",
		formattedTime.Hours,
		formattedTime.Minutes,
		formattedTime.Seconds)
	// clear start time
	startTime = 0
}

// PrintTimeout will print a timeout message for n seconds.
func PrintTimeout(duration uint32) {
	for duration >= 0 {
		fmt.Printf("\rRetrying in " + string(duration) + " seconds...")
		time.Sleep(time.Second)
		duration--
	}
	fmt.Println()
}
