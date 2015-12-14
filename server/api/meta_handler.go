package api

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/util/log"
)

func parseMetaParams(params map[string]string) *meta.Request {
	return &meta.Request{
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

func dispatchRequest(metaChan chan *meta.Response, metaReq *meta.Request) {
	// get the meta data promise
	promise := meta.GetMeta(metaReq)
	// when the meta data is ready
	promise.OnComplete(func(res interface{}) {
		metaChan <- res.(*meta.Response)
	})
}

func metaHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	// set content type response header
	w.Header().Set("Content-Type", "application/json")
	// parse params from URL
	metaReq := parseMetaParams(c.URLParams)
	// create channel to pass metadata
	metaChan := make(chan *meta.Response)
	// dispatch the request async and wait on channel
	go dispatchRequest(metaChan, metaReq)
	// wait on response
	metaRes := <-metaChan
	if metaRes.Error != nil {
		log.Warn(metaRes.Error)
		handleMetaErr(w)
		return
	}
	// send success response
	w.WriteHeader(200)
	w.Write(metaRes.Meta)
}
