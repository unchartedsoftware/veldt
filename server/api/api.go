package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/unchartedsoftware/prism/server/conf"
)

// New returns a new Goji Mux handler to process HTTP requests
func New() http.Handler {
	r := web.New()

	// Mount middleware
	r.Use(middleware.EnvInit)

	// Get conf struct
	conf := conf.GetConf()

	// Batcher websocket handler
	r.Get("/batch", batchHandler)

	// Mount tile request handler
	r.Get("/:type/:z/:x/:y", tileHandler)

	// Greedy route last
	r.Get("/*", http.FileServer(http.Dir(conf.Public)))

	return r
}
