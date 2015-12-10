package api

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/tiling"
	"github.com/unchartedsoftware/prism/util/log"
)

func parseMetaParams(params map[string]string) *tiling.TileRequest {
	return &tiling.MetaRequest{
		Endpoint: conf.Unalias(params["endpoint"]),
		Index:    conf.Unalias(params["index"]),
		Type:     conf.Unalias(params["type"]),
	}
}

func handleMetaErr(w http.ResponseWriter) {
	// send error
	w.WriteHeader(500)
	fmt.Fprint(w, `{"status": "error"}`)
}

func metaHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	// set content type response header
	w.Header().Set("Content-Type", "application/json")
	// parse params from URL
	metaReq := parseMetaParams(c.URLParams)
	// get meta data
	meta, err := tiling.GetMeta(metaReq)
	if err != nil {
		log.Warn(err)
		handleTileErr(w)
		return
	}
	// send response
	w.WriteHeader(200)
	w.Write(meta)
}
