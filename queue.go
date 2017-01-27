package veldt

import (
	"fmt"
	"runtime"
	"sync"
)

type queue struct {
	ready      chan bool
	pending    int
	mu         *sync.Mutex
	maxPending int
	maxLength  int
}

func newQueue() *queue {
	q := &queue{
		ready:      make(chan bool),
		mu:         &sync.Mutex{},
		maxPending: 32,
		maxLength:  256 * 8,
	}
	// store in intermediate here in case max is change before the following
	// loop completes
	currentMax := q.maxPending
	go func() {
		// send as many ready messages as there are expected listeners
		for i := 0; i < currentMax; i++ {
			q.ready <- true
		}
	}()
	return q
}

func (q *queue) incrementPending() error {
	q.mu.Lock()
	defer runtime.Gosched()
	defer q.mu.Unlock()
	if q.pending-q.maxPending > q.maxLength {
		return fmt.Errorf("queue has reached maximum length of %d and is no longer accepting requests",
			q.maxLength)
	}
	// increment count
	q.pending++
	return nil
}

func (q *queue) decrementPending() {
	q.mu.Lock()
	q.pending--
	q.mu.Unlock()
	runtime.Gosched()
}

func (q *queue) queueTile(req *TileRequest) ([]byte, error) {
	// increment the q.pending query count
	err := q.incrementPending()
	if err != nil {
		return nil, err
	}
	// wait until equalizer is ready
	<-q.ready
	// dispatch the query
	tile, err := req.Tile.Create(req.URI, req.Coord, req.Query)
	// decrement the q.pending count
	q.decrementPending()
	go func() {
		// inform queue that it is ready to generate another tile
		q.ready <- true
	}()
	return tile, err
}

func (q *queue) queueMeta(req *MetaRequest) ([]byte, error) {
	// increment the q.pending query count
	err := q.incrementPending()
	if err != nil {
		return nil, err
	}
	// wait until equalizer is ready
	<-q.ready
	// dispatch the query
	tile, err := req.Meta.Create(req.URI)
	// decrement the q.pending count
	q.decrementPending()
	go func() {
		// inform queue that it is ready to generate another tile
		q.ready <- true
	}()
	return tile, err
}

func (q *queue) setMaxConcurrent(max int) {
	q.mu.Lock()
	diff := max - q.maxPending
	q.maxPending = max
	q.mu.Unlock()
	if diff > 0 {
		// add ready instances to the chan
		go func() {
			// send as many ready messages as there are expected listeners
			for i := 0; i < diff; i++ {
				q.ready <- true
			}
		}()
	} else {
		// remove ready instances from the chan
		go func() {
			// send as many ready messages as there are expected listeners
			for i := diff; i < 0; i++ {
				<-q.ready
			}
		}()
	}
	runtime.Gosched()
}

func (q *queue) setQueueLength(length int) {
	q.mu.Lock()
	q.maxLength = length
	q.mu.Unlock()
	runtime.Gosched()
}
