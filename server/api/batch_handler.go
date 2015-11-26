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
	Chan chan *tiling.TileResponse
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
	dispatcher := &TileDispatcher{
		Chan: make(chan *tiling.TileResponse),
		Conn: conn,
	}
	// set the dispatcher to listen for any responses to its dispatched requests,
	// writing them out into the websocket connection as they come
	go dispatcher.ListenForResponses()
	// return the dispatcher
	return dispatcher, nil
}

// GetRequest waits on the websocket connection for tile requests.
func (t *TileDispatcher) GetRequest() (*tiling.TileRequest, error) {
	// tile request
	tileReq := &tiling.TileRequest{}
	// wait on read
	err := t.Conn.ReadJSON(&tileReq)
	if err != nil {
		return nil, err
	}
	return tileReq, nil
}

// func (t* TileDispatcher) Listen() (*tiling.TileRequest, error) {
// 	select {
// 	case err := <- t.ErrChan:
// 		return nil, err
// 	case req := <- t.ReqChan:
// 		return req, nil
// 	}
// }

// DispatchRequest takes a tile request and dispatches it to the generation package.
func (t *TileDispatcher) DispatchRequest(tileReq *tiling.TileRequest) {
	promise, err := tiling.GetTile(tileReq)
	if err != nil {
		log.Warn(err)
	}
	promise.OnComplete(func(res interface{}) {
		// cast to tile response and pass to response channel
		t.Chan <- res.(*tiling.TileResponse)
	})
}

// ListenForResponses waits on tile responses and communicates them to the client via websocket.
func (t *TileDispatcher) ListenForResponses() {
	for resp := range t.Chan {
		// write response to websocket
		t.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		t.Conn.WriteJSON(resp)
		// err := t.Conn.WriteJSON(resp)
		// if err != nil {
		// 	t.ErrChan <- err
		// }
	}
}

// Close closes the dispatchers internal channel and websocket connection.
func (t *TileDispatcher) Close() {
	// close response channel
	close(t.Chan)
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
	// begin read pump
	for {
		// wait on tile request
		tileReq, err := dispatcher.GetRequest()
		if err != nil {
			log.Debug(err)
			break
		}
		// dispatch the tile request
		dispatcher.DispatchRequest(tileReq)
	}
	// clean up dispatcher internals
	dispatcher.Close()
}
