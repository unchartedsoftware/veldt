package progress

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/unchartedsoftware/prism/util"
	"github.com/unchartedsoftware/prism/util/log"
)

var (
	startTime    time.Time
	endTime      time.Time
	currentBytes uint64
	totalBytes   uint64
	mutex        = sync.Mutex{}
	layout       = "15:04:05" // format for 3:04:05PM
)

// Time represents the hour, minutes, and seconds of a given time.
type Time struct {
	Seconds uint64
	Minutes uint64
	Hours   uint64
}

func formatTime(duration time.Duration) string {
	totalSeconds := uint64(duration.Seconds())
	totalMinutes := totalSeconds / 60
	seconds := totalSeconds % 60
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%2dh:%02dm:%02ds",
		int(math.Min(99, float64(hours))),
		minutes,
		seconds)
}

// StartProgress sets the internal epoch and the total bytes to track.
func StartProgress(bytes uint64) {
	startTime = time.Now()
	currentBytes = 0
	totalBytes = bytes
}

// EndProgress sets the end time.
func EndProgress() {
	endTime = time.Now()
}

// UpdateProgress will update and print a human readable progress message for
// a given task.
func UpdateProgress(bytes uint64) {
	mutex.Lock()
	currentBytes = currentBytes + bytes
	fCurrentBytes := float64(currentBytes)
	fTotalBytes := float64(totalBytes)
	elapsedSec := time.Since(startTime).Seconds()
	percentage := 100 * (fCurrentBytes / fTotalBytes)
	bytesPerSec := 1.0
	if elapsedSec > 0 {
		bytesPerSec = fCurrentBytes / elapsedSec
	}
	remainingBytes := fTotalBytes - fCurrentBytes
	remaining := time.Second * time.Duration(remainingBytes/bytesPerSec)
	fmtBytes := util.FormatBytes(fCurrentBytes)
	fmtBytesPerSec := util.FormatBytes(bytesPerSec)
	lineEnding := ""
	if percentage == 100 {
		// only add line ending if the progress is done
		lineEnding = "\n"
	}
	fmt.Printf("\rProcessed %+8s at %+8sps, %6.2f%% complete, estimated time remaining: %s%s",
		fmtBytes,
		fmtBytesPerSec,
		percentage,
		formatTime(remaining),
		lineEnding)
	mutex.Unlock()
}

// PrintTotalDuration prints the total duration of the processed task.
func PrintTotalDuration() {
	elapsed := endTime.Sub(startTime)
	log.Debugf("Task completed in %s",
		formatTime(elapsed))
}
