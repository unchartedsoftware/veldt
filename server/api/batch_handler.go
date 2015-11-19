/*
package api

import (
	"errors"
	"fmt"
	"strconv"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

type TileMessage struct {
	TileCoord binning.TileCoord
	Layer string
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
	defer conn.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println( err )
			break
		}

		tileMsg := &TileMessage{}
		err := json.Unmarshal( []byte(message), &tileMsg )
		if err {
			fmt.Println( err )
			continue
		}

		// get tile data and respond when it is ready
		go func() {
			elastic.GetJSONTile( tileMsg.TileCoord )
			err = c.WriteMessage(mt, message)
			if err != nil {
				fmt.Println( err )
				break
			}
		}
	}
}
*/
