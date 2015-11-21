package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/unchartedsoftware/prism/tiling"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  65536,
	WriteBufferSize: 65536,
}

func batchHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		// get tile request message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		// unmarshal the request
		tileReq := &tiling.TileRequest{}
		err = json.Unmarshal([]byte(message), &tileReq)
		if err != nil {
			fmt.Println(err)
			break
		}
		// execute tile request
		go func(tileReq *tiling.TileRequest) {
			// wait on tile response promise
			tileRes := <-tiling.GetTile(tileReq)
			message, _ = json.Marshal(&tileRes)
			conn.WriteMessage(messageType, message)
		}(tileReq)
	}
}
