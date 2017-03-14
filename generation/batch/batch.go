package batch



import (
	"sync"
	"time"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TileRequestInfo appends to a TileRequest the structures needed to get
// the request to the factory, and get the tile back from it
type tileRequestInfo struct {
	*TileRequest
	// The ID under which the factory that will fulfill this tile request was
	// registered
	factoryID     string
	// How long this tile has before its request must be made, in milliseconds
	maxWait       int64
	// An indicator of whether this tile request is ready to be fulfilled
	ready         bool
}

var (
	// Our lock to make sure we don't get race conditions in our tile request list
	mutex = sync.Mutex{}
	// Our lists of currently extant tile requests
	requests = make(map[string][]*tileRequestInfo)
	// A timer that will run all current requests when it is done
	timer *time.Timer
	// All registered tile factories
	factories = make(map[string]*TileFactory)
)

// NewBatchTile returns a tile handler for a tile request that should be
// batched with other tile requests before being processed.
//
// maxWait - the maximum time to wait before actually requesting this tile
// factory - the factory from which we 
func NewBatchTile (factoryID string, factory *TileFactory, maxWait int64) veldt.TileCtor {
	// Record our factory
	// This should only happen during pipeline construction, so concurrency
	// shouldn't be an issue
	factories[factoryID] = factory

	// And return our tile constructor function
	return func() (veldt.Tile, error) {
		batchInfof("New batched tile request")
		t := &tileRequestInfo{}
		t.resultChannel = make(chan TileResponse)
		t.maxWait = maxWait
		t.factoryID = factoryID
		t.ready   = false
		return t, nil
	}
}

// Parse records the request parameters for the factory
func (t *tileRequestInfo) Parse (params map[string]interface{}) error {
	t.parameters = params
	return nil
}

// Create waits for tile requests to come in until our waiting period is done,
// then makes any extant til requests of their tile factory, and returns the
// results to the caller
//
// uri The tile set identifier
// coord The coordinates of the desired tile
// query The filter to apply to the data when creating the tile
func (t *tileRequestInfo) Create (uri string, coords *binning.TileCoord, query veldt.Query) ([]byte, error) {
	t.uri = uri
	t.coordinates = coords
	t.query = query

	// Lock things down - we want to make sure we're listening for a response
	// from our factory before we release the lock
	mutex.Lock()
	defer mutex.Unlock()
	
	// Make sure we're all queued up
	t.enqueue()

	// And wait for a response
	response := <- t.resultChannel
	return response.tile, response.err
}

// Enqueue puts this request on the queue, and makes sure the queue is active
// Because this changes the request queue, it should only ever be called when
// our mutex lock is locked
func (t *tileRequestInfo) enqueue () {
	// Add in our request
	factoryRequests, ok := requests[t.factoryID]
	if (!ok) {
		factoryRequests = make([]*tileRequestInfo, 0)
	}
	requests[t.factoryID] = append(factoryRequests, t)

	// If no timer is running, start one
	// Ideally if one is running, we would make sure it will go off in time,
	// but that's too complex for a first pass, so we'll deal with that later.
	if nil == timer {
		timer = time.AfterFunc(time.Millisecond * time.Duration(t.maxWait), processQueue)
	}
}

// processQueue processes our request queue when the timer runs out, forwarding
// requests in bulk to the appropriate tile factory.
func processQueue () {
	// Only change the queue within a lock
	mutex.Lock()
	defer mutex.Unlock()

	for factoryID, factoryRequestInfos := range requests {
		// Get our factory
		factory := factories[factoryID]
		// Take out our meta-request info
		n := len(factoryRequestInfos)
		factoryRequests := make([]*TileRequest, n, n)
		for i := 0; i < n; i++ {
			factoryRequests[i] = factoryRequestInfos[i].TileRequest
		}

		// Call our factory, have it create tiles
		(*factory).Create(factoryRequests)
	}
	
	// All done - clear our queue, and remove our timer
	requests = make(map[string][]*tileRequestInfo)
	timer = nil
}
