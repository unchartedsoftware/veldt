package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/unchartedsoftware/prism/tiling"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 1024 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

func getTimestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func writeMessage(conn *websocket.Conn, mutex *sync.Mutex, tileReq *tiling.TileRequest) {
	// wait on tile response promise
	timestamp := getTimestamp()
	tileRes := tiling.GetTile(tileReq)
	if tileRes.Error != nil {
		fmt.Println(tileRes.Error)
	}
	fmt.Printf("Tiling request: /%s/%d/%d/%d - %dms\n",
		tileRes.Type,
		tileRes.TileCoord.Z,
		tileRes.TileCoord.X,
		tileRes.TileCoord.Y,
		getTimestamp()-timestamp)
	mutex.Lock()
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := conn.WriteJSON(tileRes)
	mutex.Unlock()
	if err != nil {
		fmt.Println(err)
	}
}

func batchHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
			break
		}
		// execute tile request
		go writeMessage(conn, mutex, tileReq)
	}
}
