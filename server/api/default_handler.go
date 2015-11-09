package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/server/conf"
)

type defaultRes struct {
	Prod bool `json:"prod"`
	Status string `json:"status"`
}

// defaultHandler handles HTTP requests for /
func defaultHandler( c web.C, w http.ResponseWriter, r *http.Request ) {
	w.Header().Set( "Content-Type", "application/json" )

	conf := conf.GetConf()
	res, err := json.Marshal(defaultRes{
		Prod: conf.Prod,
		Status: "ok",
	})

	if err != nil {
		w.WriteHeader( 500 )
		fmt.Fprint( w, `{"status": "error"}` )
		return
	}

	w.WriteHeader( 200 )
	fmt.Fprint( w, string( res ) )
}
