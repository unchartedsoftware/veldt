package queue

import (
	"fmt"
	"runtime"
	"sync"
)

// Request represents a basic request interface.
type Request interface {
	Create() ([]byte, error)
}

// Queue represents a queue for orchestating concurrent requests.
type Queue struct {
	ready      chan bool
	pending    int
	mu         *sync.Mutex
	maxPending int
	maxLength  int
}

// NewQueue instantiates and returns a new queue struct.
func NewQueue() *Queue {
	q := &Queue{
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

// Send will put the request on the queue and send it when ready.
func (q *Queue) Send(req Request) ([]byte, error) {
	// increment the q.pending query count
	err := q.incrementPending()
	if err != nil {
		return nil, err
	}
	// wait until equalizer is ready
	<-q.ready
	// dispatch the query
	res, err := req.Create()
	// decrement the q.pending count
	q.decrementPending()
	go func() {
		// inform Queue that it is ready to generate another tile
		q.ready <- true
	}()
	return res, err
}

// SetMaxConcurrent sets the maximum concurrent pending requests for the queue.
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

// SetLength sets the maximum length of the queue.
func (q *Queue) SetLength(length int) {
	q.mu.Lock()
	q.maxLength = length
	q.mu.Unlock()
	runtime.Gosched()
}

func (q *Queue) incrementPending() error {
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

func (q *Queue) decrementPending() {
	q.mu.Lock()
	q.pending--
	q.mu.Unlock()
	runtime.Gosched()
}
