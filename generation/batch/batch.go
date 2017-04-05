package batch

import (
	"time"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// tileRequestInfo appends to a TileRequest the structures needed to get
// the request to the factory, and get the tile back from it, as well as making
// the request a Tile that can be used by the pipeline.
type tileRequestInfo struct {
	TileRequest
	// The ID under which the factory that will fulfill this tile request was
	// registered
	factoryID string
	// The time at which this request was made
	time time.Time
	// A unique ID for this request
	requestID int
	// How long this tile has before its request must be made, in milliseconds
	maxWait int64
	// An indicator of whether this tile request is ready to be fulfilled
	ready bool
	// In which batch this request made
	batch int
}

var (
	// The next valid request ID
	nextRequestID = 0
)

// NewBatchTile returns a tile handler for a tile request that should be
// batched with other tile requests before being processed.
//
// maxWait - the maximum time to wait before actually requesting this tile
// factory - the factory which will actually produce our tiles
func NewBatchTile(factoryID string, factory TileFactoryCtor, maxWait int64) veldt.TileCtor {
	// Record our factory
	// This should only happen during pipeline construction, so concurrency
	// shouldn't be an issue
	factories[factoryID] = factory

	if !started {
		started = true
		go processQueue(maxWait)
	}

	// And return our tile constructor function
	return func() (veldt.Tile, error) {
		t := &tileRequestInfo{}
		t.maxWait = maxWait
		t.factoryID = factoryID
		t.ready = false
		return t, nil
	}
}

// Parse records the request parameters for the factory
func (t *tileRequestInfo) Parse(params map[string]interface{}) error {
	t.Parameters = params
	return nil
}

func (t *tileRequestInfo) getTimeDue() time.Time {
	return t.time.Add(time.Millisecond * time.Duration(t.maxWait))
}

// Create waits for tile requests to come in until our waiting period is done,
// then makes any extant tile requests of their tile factory, and returns the
// results to the caller
//
// uri The tile set identifier
// coord The coordinates of the desired tile
// query The filter to apply to the data when creating the tile
func (t *tileRequestInfo) Create(uri string, coords *binning.TileCoord, query veldt.Query) ([]byte, error) {
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
	response := <-t.ResultChannel
	batchInfof("Request %d for tile set %s, factory id %s, tile %v fulfilled at %v", t.requestID, t.URI, t.factoryID, coords, time.Now())
	// Done with this channel; close it.
	close(t.ResultChannel)

	return response.Tile, response.Err
}

// Enqueue puts this request on the queue, and makes sure the queue is active
// Because this changes the request queue, it should only ever be called when
// our mutex lock is locked
func (t *tileRequestInfo) enqueue() {
	// Only alter the queue under lock
	mutex.Lock()
	defer mutex.Unlock()

	// Add in our request
	factoryRequests, ok := requests[t.factoryID]
	if !ok {
		factoryRequests = make([]*tileRequestInfo, 0)
	}
	requests[t.factoryID] = append(factoryRequests, t)
}
