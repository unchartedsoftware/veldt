package api

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/util/log"
)

func parseMetaParams(params map[string]string) *meta.MetaRequest {
	return &meta.MetaRequest{
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
	// create channel to pass metadata
	metaChan := make(chan *meta.MetaResponse)
	// get the meta data promise
	promise := meta.GetMeta(metaReq)
	// when the meta data is ready
	promise.OnComplete(func(res interface{}) {
		metaChan <- res.(*meta.MetaResponse)
	})
	// wait on response
	metaRes := <-metaChan
	if metaRes.Error != nil {
		log.Warn(metaRes.Error)
		handleTileErr(w)
		return
	}
	// send success response
	w.WriteHeader(200)
	w.Write(metaRes.Meta)
}
