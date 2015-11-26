package api

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"

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
	RespChan chan *tiling.TileResponse
	ErrChan chan error
	Conn *websocket.Conn
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
		RespChan: make(chan *tiling.TileResponse),
		ErrChan: make(chan error),
		Conn: conn,
	}, nil
}

// Listen waits on both tile request and responses and handles each until the websocket connection dies.
func (t* TileDispatcher) Listen() error {
	go t.listenForRequests()
	go t.listenForResponses()
	return <- t.ErrChan
}

// ListenForResponses waits on tile responses and communicates them to the client via websocket.
func (t *TileDispatcher) listenForResponses() {
	for resp := range t.RespChan {
		// write response to websocket
		t.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		err := t.Conn.WriteJSON(resp)
		if err != nil {
			t.ErrChan <- err
			break
		}
	}
}

// DispatchRequest takes a tile request and dispatches it to the generation package.
func (t *TileDispatcher) dispatchRequest(tileReq *tiling.TileRequest) {
	promise, err := tiling.GetTile(tileReq)
	if err != nil {
		log.Warn(err)
	}
	promise.OnComplete(func(res interface{}) {
		// cast to tile response and pass to response channel
		t.RespChan <- res.(*tiling.TileResponse)
	})
}

// GetRequest waits on the websocket connection for tile requests.
func (t *TileDispatcher) getRequest() (*tiling.TileRequest, error) {
	// tile request
	tileReq := &tiling.TileRequest{}
	// wait on read
	err := t.Conn.ReadJSON(&tileReq)
	if err != nil {
		return nil, err
	}
	return tileReq, nil
}

// ListenForResponses waits on tile responses and communicates them to the client via websocket.
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

// Close closes the dispatchers internal channel and websocket connection.
func (t *TileDispatcher) Close() {
	// close response channel
	close(t.RespChan)
	// close error channel
	close(t.ErrChan)
	// close websocket connection
	t.Conn.Close()
}

func batchHandler(w http.ResponseWriter, r *http.Request) {
	// create dispatcher
	dispatcher, err := NewTileDispatcher(w, r)
	if err != nil {
		log.Warn(err)
		return
	}
	err = dispatcher.Listen()
	if err != nil {
		log.Debug(err)
	}
	// clean up dispatcher internals
	dispatcher.Close()
}
