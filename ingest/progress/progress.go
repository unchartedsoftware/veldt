package progress

import (
    "fmt"
    "time"
)

type Time struct {
    Seconds uint64
    Minutes uint64
    Hours uint64
}

func formatTime( totalSeconds uint64 ) Time {
    totalMinutes := totalSeconds / 60
    seconds := totalSeconds % 60
    hours := totalMinutes / 60
    minutes := totalMinutes % 60
    return Time{
        Seconds: seconds,
        Minutes: minutes,
        Hours: hours,
    }
}

func getTimestamp() uint64 {
    return uint64( time.Now().Unix() )
}

var startTime uint64 = 0

func PrintProgress( totalBytes int64, bytes int64 ) {
    if startTime == 0 {
        startTime = getTimestamp()
    }
    elapsed := getTimestamp() - startTime
    percentComplete := 100 * ( float64( bytes ) / float64( totalBytes ) )
    bytesPerSecond := 1.0
    if elapsed > 0 {
        bytesPerSecond = float64( bytes ) / float64( elapsed )
    }
    estimatedSecondsRemaining := ( float64( totalBytes ) - float64( bytes ) ) / bytesPerSecond
    formattedTime := formatTime( uint64( estimatedSecondsRemaining ) )
    fmt.Printf( "\rProcessed %d bytes at %f Bps, %f%% complete, estimated time remaining %d:%02d:%02d",
        bytes,
        bytesPerSecond,
        percentComplete,
        formattedTime.Hours,
        formattedTime.Minutes,
        formattedTime.Seconds )
}

func PrintTotalDuration() {
    // finished succesfully
    formattedTime := formatTime( getTimestamp() - startTime )
    fmt.Printf( "\nTask completed in %d:%02d:%02d\n",
        formattedTime.Hours,
        formattedTime.Minutes,
        formattedTime.Seconds )
}

func PrintTimeout( duration uint32 ) {
    for duration >= 0 {
        fmt.Printf( "\rRetrying in " + string( duration ) + " seconds..." )
        time.Sleep( time.Second )
        duration -= 1
    }
    fmt.Println()
}
