package throttle

import (
	"fmt"
	"runtime"
	"sync"

	"gopkg.in/olivere/elastic.v3"
)

const (
	// max number of concurrent requests
	maxPendingRequests = 8
	// max number of queries to queue until returning errors
	maxQueueLength = 64
)

var (
	ready   = make(chan bool)
	pending = 0
	mutex   = sync.Mutex{}
)

func init() {
	go func() {
		// send as many ready messages as there are expected listeners
		for i := 0; i < maxPendingRequests; i++ {
			ready <- true
		}
	}()
}

func incrementPending() error {
	mutex.Lock()
	defer runtime.Gosched()
	defer mutex.Unlock()
	if pending-maxPendingRequests > maxQueueLength {
		return fmt.Errorf("Elasticsearch queue has reached maximum length of %d and is no longer accepting requests", maxQueueLength)
	}
	// increment count
	pending++
	return nil
}

func decrementPending() {
	mutex.Lock()
	pending--
	mutex.Unlock()
	runtime.Gosched()
}

// Send dispatches an elastic search query, limiting the number of concurrent
// requests.
func Send(req *elastic.SearchService) (*elastic.SearchResult, error) {
	// increment the pending query count
	err := incrementPending()
	if err != nil {
		return nil, err
	}
	// wait until equalizer is ready
	<-ready
	// dispatch the query
	res, err := req.Do()
	// decrement the pending count
	decrementPending()
	go func() {
		// inform eq that it is ready to pick up another query
		ready <- true
	}()
	return res, err
}
