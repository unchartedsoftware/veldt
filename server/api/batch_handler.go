package api

import (
	"net/http"
	"sync"
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

func writeMessage(conn *websocket.Conn, mutex *sync.Mutex, tileReq *tiling.TileRequest) {
	p, err := tiling.GetTile(tileReq)
	if err != nil {
		log.Warn(err)
	}
	p.OnComplete(func(res interface{}) {
		// cast to tile response
		tileRes := res.(*tiling.TileResponse)
		// writes are not thread-safe
		mutex.Lock()
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		err = conn.WriteJSON(tileRes)
		mutex.Unlock()
		if err != nil {
			log.Warn(err)
		}
	})
}

func batchHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(err)
		return
	}
	defer conn.Close()
	conn.SetReadLimit(maxMessageSize)
	mutex := &sync.Mutex{}
	for {
		// tile request
		tileReq := &tiling.TileRequest{}
		// wait on read
		err := conn.ReadJSON(&tileReq)
		if err != nil {
			log.Debug(err)
			break
		}
		// execute tile request
		go writeMessage(conn, mutex, tileReq)
	}
}
