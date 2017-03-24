package batch



import (
	"fmt"
	"sync"
	"time"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TileRequestInfo appends to a TileRequest the structures needed to get
// the request to the factory, and get the tile back from it
type tileRequestInfo struct {
	TileRequest
	// The ID under which the factory that will fulfill this tile request was
	// registered
	factoryID     string
	// The time at which this request was made
	time          time.Time
	// A unique ID for this request
	requestID     int
	// How long this tile has before its request must be made, in milliseconds
	maxWait       int64
	// An indicator of whether this tile request is ready to be fulfilled
	ready         bool
	// In which batch this request made
	batch         int
}

var (
	// Our lock to make sure we don't get race conditions in our tile request list
	mutex = sync.Mutex{}
	// Our lists of currently extant tile requests
	requests = make(map[string][]*tileRequestInfo)
	// All registered tile factories
	factories = make(map[string]TileFactoryCtor)
	// The current queue batch
	queueBatch = 0
	// The next valid request ID
	nextRequestID = 0
	// Is our event processing loop started?
	started = false
)

// NewBatchTile returns a tile handler for a tile request that should be
// batched with other tile requests before being processed.
//
// maxWait - the maximum time to wait before actually requesting this tile
// factory - the factory from which we 
func NewBatchTile (factoryID string, factory TileFactoryCtor, maxWait int64) veldt.TileCtor {
	// Record our factory
	// This should only happen during pipeline construction, so concurrency
	// shouldn't be an issue
	factories[factoryID] = factory

	if (!started) {
		started = true
		go processQueue(maxWait)
	}

	// And return our tile constructor function
	return func() (veldt.Tile, error) {
		t := &tileRequestInfo{}
		t.maxWait = maxWait
		t.factoryID = factoryID
		t.ready   = false
		return t, nil
	}
}

// Parse records the request parameters for the factory
func (t *tileRequestInfo) Parse (params map[string]interface{}) error {
	t.Parameters = params
	return nil
}

func (t *tileRequestInfo) getTimeDue () time.Time {
	return t.time.Add(time.Millisecond * time.Duration(t.maxWait))
}

// Create waits for tile requests to come in until our waiting period is done,
// then makes any extant tile requests of their tile factory, and returns the
// results to the caller
//
// uri The tile set identifier
// coord The coordinates of the desired tile
// query The filter to apply to the data when creating the tile
func (t *tileRequestInfo) Create (uri string, coords *binning.TileCoord, query veldt.Query) ([]byte, error) {
	t.URI = uri
	t.Coordinates = coords
	t.Query = query
	t.ResultChannel = make(chan TileResponse, 1)	
	t.time = time.Now()
	mutex.Lock()
	t.requestID = nextRequestID
	nextRequestID = nextRequestID + 1
	mutex.Unlock()


	batchDebugf("Queueing up request for tile set %s, factory id %s, tile %v", uri, t.factoryID, coords)
	t.enqueue()

	batchInfof("Request %d for tile set %s, factory id %s, tile %v enqueued at %v", t.requestID, t.URI, t.factoryID, coords, t.time)
	response := <- t.ResultChannel
	batchInfof("Request %d for tile set %s, factory id %s, tile %v fulfilled at %v", t.requestID, t.URI, t.factoryID, coords, time.Now())
	return response.Tile, response.Err
}

// Enqueue puts this request on the queue, and makes sure the queue is active
// Because this changes the request queue, it should only ever be called when
// our mutex lock is locked
func (t *tileRequestInfo) enqueue () {
	// Only alter the queue under lock
	mutex.Lock()
	defer mutex.Unlock()

	// Add in our request
	factoryRequests, ok := requests[t.factoryID]
	if (!ok) {
		factoryRequests = make([]*tileRequestInfo, 0)
	}
	requests[t.factoryID] = append(factoryRequests, t)
}

func processQueue (waitTime int64) {
	for {
		batchInfof("Checking for queued tiles at %v", time.Now())
		if (isDue()) {
			batch, batchRequests := dequeueRequests()
			numRequests := 0
			for _, factoryRequests := range batchRequests {
				numRequests = numRequests + len(factoryRequests)
			}
			batchInfof("Beginning processing of batch %d, with %d requests in %d factories, at %v", batch, numRequests, len(batchRequests), time.Now())
			for factoryID, factoryRequests := range batchRequests {
				processFactoryRequests(batch, factoryID, factoryRequests)
			}
			batchInfof("Done processing of batch %d at %v", batch, time.Now())
		}
		time.Sleep(time.Millisecond * time.Duration(waitTime))
	}
}

// dequeueRequests takes requests off of the queue in preparation for sending
// them to their respective factories
func dequeueRequests () (int, map[string][]*tileRequestInfo) {
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
func processFactoryRequests (batch int, factoryID string, factoryRequests []*tileRequestInfo) {
	batchDebugf("Processing %d requests for factory %s", len(factoryRequests), factoryID)

	// Get our factory
	ctor, ok := factories[factoryID]
	if !ok {
		err := fmt.Errorf("Unrecognized tile factory '%s'", factoryID)
		batchWarnf(err.Error())
		sendError(err, factoryRequests)
	} else {
		factory, err := ctor()
		if nil != err {
			err := fmt.Errorf("Error constructing factory %s: %v", factoryID, err)
			batchWarnf(err.Error())
			sendError(err, factoryRequests)
		} else {
			batchDebugf("Factory obtained.")

			// Take out our meta-request info, leaving just the simple request info
			// for the factory
			n := len(factoryRequests)
			simpleRequests := make([]*TileRequest, n, n)
			for i := 0; i < n; i++ {
				batchDebugf("request: factory=%s, batch=%d, uri=%s, coords=%v",
					factoryID, factoryRequests[i].batch,
					factoryRequests[i].URI, factoryRequests[i].Coordinates)
				simpleRequests[i] = &factoryRequests[i].TileRequest
			}

			// Call our factory, have it create tiles
			batchDebugf("\tCalling factory %s to create tiles", factoryID)
			factory.CreateTiles(simpleRequests)
		}
	}
}

func sendError (err error, requestInfos []*tileRequestInfo) {
	response := TileResponse{nil, err}
	for _, requestInfo := range requestInfos {
		requestInfo.ResultChannel <- response
	}
}

func isDue () bool {
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
