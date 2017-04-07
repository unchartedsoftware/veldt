package batch

import (
	"fmt"
	"sync"
	"time"
)

var (
	// Our lock to make sure we don't get race conditions in our tile request list
	mutex = sync.Mutex{}
	// Our lists of currently extant tile requests
	requests = make(map[string][]*tileRequestInfo)
	// All registered tile factories
	factories = make(map[string]TileFactoryCtor)
	// The current queue batch
	queueBatch = 0
	// Is our event processing loop started?
	started = false
)

func processQueue(waitTime int64) {
	for {
		Infof("Checking for queued tiles at %v", time.Now())
		if isDue() {
			batch, batchRequests := dequeueRequests()
			numRequests := 0
			for _, factoryRequests := range batchRequests {
				numRequests = numRequests + len(factoryRequests)
			}
			Infof("Beginning processing of batch %d, with %d requests in %d factories, at %v", batch, numRequests, len(batchRequests), time.Now())
			for factoryID, factoryRequests := range batchRequests {
				processFactoryRequests(batch, factoryID, factoryRequests)
			}
			Infof("Done processing of batch %d at %v", batch, time.Now())
		}
		time.Sleep(time.Millisecond * time.Duration(waitTime))
	}
}

// dequeueRequests takes requests off of the queue in preparation for sending
// them to their respective factories
func dequeueRequests() (int, map[string][]*tileRequestInfo) {
	// Only change the queue within a lock
	mutex.Lock()
	defer mutex.Unlock()

	// Update our queue batch
	queueBatch = queueBatch + 1

	// Retrieve our current set of requests, setting up a new collector for
	// subsequent requests.
	current := requests
	requests = make(map[string][]*tileRequestInfo)

	// Mark the batch number on all current requests
	for _, requestSet := range current {
		for _, req := range requestSet {
			req.batch = queueBatch
		}
	}

	return queueBatch, current
}

// processFactoryRequess process all the requests from a given batch for a
// given factory
func processFactoryRequests(batch int, factoryID string, factoryRequests []*tileRequestInfo) {
	Debugf("Processing %d requests for factory %s", len(factoryRequests), factoryID)

	// Get our factory
	ctor, ok := factories[factoryID]
	if !ok {
		err := fmt.Errorf("unrecognized tile factory '%s'", factoryID)
		Warnf(err.Error())
		sendError(err, factoryRequests)
		return
	}

	factory, err := ctor()
	if nil != err {
		err := fmt.Errorf("error constructing factory %s: %v", factoryID, err)
		Warnf(err.Error())
		sendError(err, factoryRequests)
		return
	}

	Debugf("Factory obtained.")
	// Take out our meta-request info, leaving just the simple request info
	// for the factory
	n := len(factoryRequests)
	simpleRequests := make([]*TileRequest, n, n)
	for i := 0; i < n; i++ {
		Debugf("request: factory=%s, batch=%d, uri=%s, coords=%v",
			factoryID, factoryRequests[i].batch,
			factoryRequests[i].URI, factoryRequests[i].Coordinates)
		simpleRequests[i] = &factoryRequests[i].TileRequest
	}

	// Call our factory, have it create tiles
	Debugf("\tCalling factory %s to create tiles", factoryID)
	factory.CreateTiles(simpleRequests)
}

func sendError(err error, requestInfos []*tileRequestInfo) {
	response := TileResponse{nil, err}
	for _, requestInfo := range requestInfos {
		requestInfo.ResultChannel <- response
	}
}

func isDue() bool {
	// We are reading our requests lists, therefor must operatate inside a lock
	mutex.Lock()
	defer mutex.Unlock()

	for _, factoryRequests := range requests {
		for _, request := range factoryRequests {
			if time.Now().After(request.getTimeDue()) {
				return true
			}
		}
	}
	return false
}
