package info

import (
	"github.com/unchartedsoftware/prism/ingest/es"
	"github.com/unchartedsoftware/prism/ingest/progress"
)

// Worker represents a designated worker function to batch in a pool.
type Worker func(IngestFile, *es.Equalizer) error

// Pool represents a single goroutine pool for batching workers.
type Pool struct {
	FileChan chan IngestFile
	ErrChan  chan error
	KillChan chan bool
	Size     int
}

// NewPool returns a new pool object with the given worker size
func NewPool(size int) *Pool {
	return &Pool{
		FileChan: make(chan IngestFile),
		ErrChan:  make(chan error),
		KillChan: make(chan bool),
		Size:     size,
	}
}

func workerWrapper(p *Pool, eq *es.Equalizer, worker Worker, ingestInfo *IngestInfo) {
	// tell the pool that this worker is ready
	p.ErrChan <- nil
	// begin worker loop
	for {
		select {
		case file := <-p.FileChan:
			// do work
			err := worker(file, eq)
			// broadcast work response to pool, if nil worker is ready for more
			// work, if not, then shut down the pool
			p.ErrChan <- err
			// if no error, print current progress
			if err == nil {
				// Update and print current progress
				progress.UpdateProgress(file.Size)
			}
		case <-p.KillChan:
			// kill worker
			return
		}
	}
}

// Close safely closes the pool
func (p *Pool) Close() {
	// workers will currently be blocked trying to send a ready/error message
	// to the pool. We need to absorb those messages now so that they will be
	// able to process the kill signals.
	for i := 0; i < p.Size; i++ {
		go func() {
			<-p.ErrChan
		}()
	}
	// send a kill message to all workers
	for i := 0; i < p.Size; i++ {
		p.KillChan <- true
	}
	// workers are all dead now, it is safe to close the channels
	close(p.FileChan)
	close(p.KillChan)
	close(p.ErrChan)
}

// Execute launches a batch of ingest workers with the provided ingest information.
func (p *Pool) Execute(worker Worker, ingestInfo *IngestInfo) error {
	// create equalizer of same size
	eq := es.NewEqualizer(p.Size)
	eq.Listen()

	// close the equalizer AFTER closing the pool, otherwise the equalizer may
	// send on a closed channel
	defer eq.Close()

	// for each worker in pool
	for i := 0; i < p.Size; i++ {
		// dispatch the workers, they will wait until the input channel is closed
		go workerWrapper(p, eq, worker, ingestInfo)
	}

	// start progress tracking
	progress.StartProgress(ingestInfo.NumTotalBytes)

	// process all files by spreading them to free workers, this blocks until
	// a worker is available, or exits if there is an error
	for _, file := range ingestInfo.Files {
		err := <-p.ErrChan
		if err != nil {
			// error has occured, close the pool and return error
			p.Close()
			return err
		}
		// if not, continue passing files to workers
		p.FileChan <- file
	}

	// when work is done safely close the pool
	p.Close()

	// end progress tracking, and print summary
	progress.EndProgress()
	progress.PrintTotalDuration()
	return nil
}
