package tile

import (
	"fmt"
	"runtime"
	"sync"
)

type Queue struct {
	ready chan bool
	q.pending int
	mu sync.Mutex
	maxPending int
	maxLength int
}

func NewQueue() *Queue {
	queue := &Queue{
		ready: make(chan bool)
		mu: sync.Mutex{},
		maxPending: 32,
		maxLength: 256 * 8,
	}
	// store in intermediate here in case max is change before the following
	// loop completes
	currentMax := queue.maxPending
	go func() {
		// send as many ready messages as there are expected listeners
		for i := 0; i < currentMax; i++ {
			queue.ready <- true
		}
	}()
	return queue
}

func (q *Queue) incrementPending() error {
	q.mu.Lock()
	defer runtime.Gosched()
	defer q.mu.Unlock()
	if q.pending-q.maxPending > q.maxLength {
		return fmt.Errorf("Queue has reached maximum length of %d and is no longer accepting requests",
			q.maxLength)
	}
	// increment count
	q.pending++
	return nil
}

func (q *Queue) decrementPending() {
	q.mu.Lock()
	q.pending--
	q.mu.Unlock()
	runtime.Gosched()
}

func (q *Queue) QueueTile(req *prism.TileRequest) ([]byte, error) {
	// increment the q.pending query count
	err := incrementPending()
	if err != nil {
		return nil, err
	}
	// wait until equalizer is ready
	<-ready
	// dispatch the query
	tile, err := req.Tile.Create(req.URI, req.Coord)
	// decrement the q.pending count
	decrementPending()
	go func() {
		// inform queue that it is ready to generate another tile
		ready <- true
	}()
	return tile, err
}

func (q *Queue) QueueMeta(req *prism.MetaRequest) ([]byte, error) {
	// increment the q.pending query count
	err := incrementPending()
	if err != nil {
		return nil, err
	}
	// wait until equalizer is ready
	<-ready
	// dispatch the query
	tile, err := req.Meta.Create(req.URI)
	// decrement the q.pending count
	decrementPending()
	go func() {
		// inform queue that it is ready to generate another tile
		ready <- true
	}()
	return tile, err
}

// SetMaxConcurrent sets the maximum concurrent tile requests allowed.
func (q *Queue) SetMaxConcurrent(max int) {
	q.mu.Lock()
	diff := max - q.maxPending
	q.maxPending = max
	q.mu.Unlock()
	if diff > 0 {
		// add ready instances to the chan
		go func() {
			// send as many ready messages as there are expected listeners
			for i := 0; i < diff; i++ {
				ready <- true
			}
		}()
	} else {
		// remove ready instances from the chan
		go func() {
			// send as many ready messages as there are expected listeners
			for i := diff; i < 0; i++ {
				<-ready
			}
		}()
	}
	runtime.Gosched()
}

// SetQueueLength sets the queue length for tiles to hold in the queue.
func (q *Queue) SetQueueLength(length int) {
	q.mu.Lock()
	q.maxLength = length
	q.mu.Unlock()
	runtime.Gosched()
}
