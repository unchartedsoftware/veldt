package pool

import (
	"sync"

	"github.com/unchartedsoftware/prism/ingest/info"
	"github.com/unchartedsoftware/prism/ingest/progress"
)

// Worker represents a designated worker function to batch in a pool.
type Worker func(info.IngestFile)

// Pool represents a single goroutine pool for batching workers.
type Pool struct {
	Channel   chan info.IngestFile
	WaitGroup *sync.WaitGroup
	Size      int
}

// New returns a new pool object with the given worker size
func New(size int) Pool {
	return Pool{
		Channel:   make(chan info.IngestFile),
		WaitGroup: new(sync.WaitGroup),
		Size:      size,
	}
}

// Track how many bytes of data has been processed
var numProcessedBytes = uint64(0)

func workerWrapper(fileChan chan info.IngestFile, waitGroup *sync.WaitGroup, worker Worker, ingestInfo *info.IngestInfo) {
	// Decrease internal counter for wait-group as soon as goroutine finishes
	defer waitGroup.Done()
	for file := range fileChan {
		// Print current progress
		worker(file)
		// Increment processed bytes
		numProcessedBytes += file.Size
		// Print current progress
		progress.PrintProgress(ingestInfo.NumTotalBytes, numProcessedBytes)
	}
}

// Execute launches a batch of ingest workers with the provided ingest information.
func (p *Pool) Execute(worker Worker, ingestInfo *info.IngestInfo) {
	// Adding routines to workgroup and running then
	for i := 0; i < p.Size; i++ {
		p.WaitGroup.Add(1)
		go workerWrapper(p.Channel, p.WaitGroup, worker, ingestInfo)
	}
	// Processing all links by spreading them to `free` goroutines
	for _, file := range ingestInfo.Files {
		p.Channel <- file
	}
	// Closing channel (waiting in goroutines won't continue any more)
	close(p.Channel)
	// Waiting for all goroutines to finish (otherwise they die as main routine dies)
	p.WaitGroup.Wait()
}
