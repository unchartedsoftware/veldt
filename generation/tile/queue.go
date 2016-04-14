package tile

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	ready   = make(chan bool)
	pending = 0
	mutex   = sync.Mutex{}
	// max number of concurrent requests
	maxPendingRequests = 32
	// max number of queries to queue until returning errors
	maxQueueLength = 256 * 8
)

func init() {
	// store in intermediate here in case max is change before the following
	// loop completes
	currentMax := maxPendingRequests
	go func() {
		// send as many ready messages as there are expected listeners
		for i := 0; i < currentMax; i++ {
			ready <- true
		}
	}()
}

func incrementPending() error {
	mutex.Lock()
	defer runtime.Gosched()
	defer mutex.Unlock()
	if pending-maxPendingRequests > maxQueueLength {
		return fmt.Errorf("Tile generation queue has reached maximum length of %d and is no longer accepting requests", maxQueueLength)
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

func queue(gen Generator) ([]byte, error) {
	// increment the pending query count
	err := incrementPending()
	if err != nil {
		return nil, err
	}
	// wait until equalizer is ready
	<-ready
	// dispatch the query
	tile, err := gen.GetTile()
	// decrement the pending count
	decrementPending()
	go func() {
		// inform queue that it is ready to generate another tile
		ready <- true
	}()
	return tile, err
}

// SetMaxConcurrent sets the maximum concurrent tile requests allowed.
func SetMaxConcurrent(max int) {
	mutex.Lock()
	diff := max - maxPendingRequests
	maxPendingRequests = max
	mutex.Unlock()
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
func SetQueueLength(length int) {
	mutex.Lock()
	maxQueueLength = length
	mutex.Unlock()
	runtime.Gosched()
}
