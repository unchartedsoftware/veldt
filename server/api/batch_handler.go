package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/tiling"
	"github.com/unchartedsoftware/prism/util/log"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 256 * 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

// TileDispatcher represents a single clients tile dispatcher.
type TileDispatcher struct {
	RespChan  chan *tiling.TileResponse
	ErrChan   chan error
	WaitGroup *sync.WaitGroup
	Conn      *websocket.Conn
}

// NewTileDispatcher returns a pointer to a new tile dispatcher object.
func NewTileDispatcher(w http.ResponseWriter, r *http.Request) (*TileDispatcher, error) {
	// open a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	// set the message read limit
	conn.SetReadLimit(maxMessageSize)
	return &TileDispatcher{
		RespChan:  make(chan *tiling.TileResponse),
		ErrChan:   make(chan error),
		WaitGroup: new(sync.WaitGroup),
		Conn:      conn,
	}, nil
}

// Listen waits on both tile request and responses and handles each until the websocket connection dies.
func (t *TileDispatcher) ListenAndRespond() error {
	go t.listenForRequests()
	go t.listenForResponses()
	return <-t.ErrChan
}

// Close closes the dispatchers internal channels and websocket connection.
func (t *TileDispatcher) Close() {
	// wait to ensure that no more responses are pending
	t.WaitGroup.Wait()
	// close dispatcher channels
	close(t.RespChan)
	close(t.ErrChan)
	// close websocket connection
	t.Conn.Close()
}

func (t *TileDispatcher) listenForResponses() {
	for tileRes := range t.RespChan {
		// log error if there is one
		if tileRes.Error != nil {
			log.Warn(tileRes.Error)
		}
		// alias endpoint, index, and type
		tileRes.Endpoint = conf.Alias(tileRes.Endpoint)
		tileRes.Index = conf.Alias(tileRes.Index)
		tileRes.Type = conf.Alias(tileRes.Type)
		// write response to websocket
		t.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		err := t.Conn.WriteJSON(tileRes)
		if err != nil {
			t.ErrChan <- err
			break
		}
	}
}

func (t *TileDispatcher) dispatchRequest(tileReq *tiling.TileRequest) {
	// increment pending response wait group to ensure we don't send on
	// a closed channel
	t.WaitGroup.Add(1)
	// get the tile promise
	promise := tiling.GetTile(tileReq)
	// when the tile is ready
	promise.OnComplete(func(res interface{}) {
		// cast to tile response and pass to response channel
		t.RespChan <- res.(*tiling.TileResponse)
		// decrement the pending response wait group
		t.WaitGroup.Done()
	})
}

func (t *TileDispatcher) getRequest() (*tiling.TileRequest, error) {
	// tile request
	tileReq := &tiling.TileRequest{}
	// wait on read
	err := t.Conn.ReadJSON(&tileReq)
	if err != nil {
		return nil, err
	}
	// unalias endpoint, index, and type
	tileReq.Endpoint = conf.Unalias(tileReq.Endpoint)
	tileReq.Index = conf.Unalias(tileReq.Index)
	tileReq.Type = conf.Unalias(tileReq.Type)
	return tileReq, nil
}

func (t *TileDispatcher) listenForRequests() {
	for {
		// wait on tile request
		tileReq, err := t.getRequest()
		if err != nil {
			t.ErrChan <- err
			break
		}
		// dispatch the request
		go t.dispatchRequest(tileReq)
	}
}

func batchHandler(w http.ResponseWriter, r *http.Request) {
	// create dispatcher
	dispatcher, err := NewTileDispatcher(w, r)
	if err != nil {
		log.Warn(err)
		return
	}
	err = dispatcher.ListenAndRespond()
	if err != nil {
		log.Debug(err)
	}
	// clean up dispatcher internals
	dispatcher.Close()
}
