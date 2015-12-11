package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/log"
)

func parseTileParams(params map[string]string) (*tile.TileRequest, error) {
	x, ex := strconv.ParseUint(params["x"], 10, 32)
	y, ey := strconv.ParseUint(params["y"], 10, 32)
	z, ez := strconv.ParseUint(params["z"], 10, 32)
	if ex == nil || ey == nil || ez == nil {
		return &tile.TileRequest{
			TileCoord: binning.TileCoord{
				X: uint32(x),
				Y: uint32(y),
				Z: uint32(z),
			},
			Endpoint: conf.Unalias(params["endpoint"]),
			Index:    conf.Unalias(params["index"]),
			Type:     conf.Unalias(params["type"]),
		}, nil
	}
	return nil, errors.New("Unable to parse tile coordinate from URL")
}

func handleTileErr(w http.ResponseWriter) {
	// send error
	w.WriteHeader(500)
	fmt.Fprint(w, `{"status": "error"}`)
}

func tileHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	// set content type response header
	w.Header().Set("Content-Type", "application/json")
	// parse tile coord from URL
	tileReq, parseErr := parseTileParams(c.URLParams)
	if parseErr != nil {
		log.Warn(parseErr)
		handleTileErr(w)
		return
	}
	// get tile hash
	tileHash := tile.GetTileHash(tileReq)
	// get tile data from store
	tileData, tileErr := store.Get(tileHash)
	if tileData == nil || tileErr != nil {
		log.Warn(tileErr)
		handleTileErr(w)
		return
	}
	// send response
	w.WriteHeader(200)
	w.Write(tileData)
}
