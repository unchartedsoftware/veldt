package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

// New returns a new Goji Mux handler to process HTTP requests
func New() http.Handler {
	r := web.New()

	// Mount middleware
	r.Use( middleware.EnvInit )

	// Mount request handlers
	r.Get( "/", defaultHandler )

	return r
}
