package batch

import (
	"net/http"

	"github.com/unchartedsoftware/prism/log"
)

const (
	// Route represents the HTTP route for the resource.
	Route = "/batch"
)

// Handler represents the HTTP route response handler.
func Handler(w http.ResponseWriter, r *http.Request) {
	// create dispatcher
	dispatcher, err := NewTileDispatcher(w, r)
	if err != nil {
		log.Warn(err)
		return
	}
	// listen for requests and respond
	err = dispatcher.ListenAndRespond()
	if err != nil {
		log.Debug(err)
	}
	// clean up dispatcher internals
	dispatcher.Close()
}
