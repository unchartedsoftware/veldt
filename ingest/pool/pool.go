package pool

import (
    "os"
    "sync"

    "github.com/unchartedsoftware/prism/ingest/info"
    "github.com/unchartedsoftware/prism/ingest/progress"
)

type Worker func( file os.FileInfo )

type Pool struct {
    Channel chan os.FileInfo
    WaitGroup *sync.WaitGroup
    Size int
}

// New returns a new pool object with the given worker size
func New( size int ) Pool {
    return Pool{
        Channel: make( chan os.FileInfo ),
        WaitGroup: new( sync.WaitGroup ),
        Size: size,
    }
}

// Track how many bytes of data has been processed
var numIngestedBytes = int64( 0 )

func workerWrapper( linkChan chan os.FileInfo, waitGroup *sync.WaitGroup, worker Worker, ingestInfo *info.IngestInfo ) {
    // Decrease internal counter for wait-group as soon as goroutine finishes
    defer waitGroup.Done()
	for file := range linkChan {
        // Print current progress
        worker( file )
        // Increment ingested bytes
        numIngestedBytes += file.Size()
        // Print current progress
        progress.PrintProgress( ingestInfo.NumTotalBytes, numIngestedBytes )
	}
}

func (p *Pool) Execute( worker Worker, ingestInfo *info.IngestInfo ) {
    // Adding routines to workgroup and running then
    for i := 0; i<p.Size; i++ {
        p.WaitGroup.Add( 1 )
        go workerWrapper( p.Channel, p.WaitGroup, worker, ingestInfo )
    }
    // Processing all links by spreading them to `free` goroutines
    for _, file := range ingestInfo.Files {
        p.Channel <- file
    }
    // Closing channel (waiting in goroutines won't continue any more)
    close( p.Channel )
    // Waiting for all goroutines to finish (otherwise they die as main routine dies)
    p.WaitGroup.Wait()
}
