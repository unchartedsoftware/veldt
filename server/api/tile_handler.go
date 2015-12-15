package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/util/log"
)

func parseParams(params url.Values) map[string]interface{} {
	p := make(map[string]interface{})
	for k, _ := range params {
		p[k] = params.Get(k)
	}
	return p
}

// func parseParams(p map[string]string) map[string]interface{} {
// 	params := make(map[string]interface{})
// 	for path, val := range p {
// 		json.Set(params, val, strings.Split(path, "."))
// 	}
// 	return params
// }

func parseTileParams(params map[string]string, queryParams url.Values) (*tile.Request, error) {
	x, ex := strconv.ParseUint(params["x"], 10, 32)
	y, ey := strconv.ParseUint(params["y"], 10, 32)
	z, ez := strconv.ParseUint(params["z"], 10, 32)
	if ex == nil || ey == nil || ez == nil {
		return &tile.Request{
			TileCoord: binning.TileCoord{
				X: uint32(x),
				Y: uint32(y),
				Z: uint32(z),
			},
			Endpoint: conf.Unalias(params["endpoint"]),
			Index:    conf.Unalias(params["index"]),
			Type:     conf.Unalias(params["type"]),
			Params:   parseParams(queryParams),
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
	tileReq, parseErr := parseTileParams(c.URLParams, r.URL.Query())
	if parseErr != nil {
		log.Warn(parseErr)
		handleTileErr(w)
		return
	}
	// get tile hash
	tileData, tileErr := tile.GetTileFromStore(tileReq)
	if tileErr != nil {
		log.Warn(tileErr)
		handleTileErr(w)
		return
	}
	// send response
	w.WriteHeader(200)
	w.Write(tileData)
}
