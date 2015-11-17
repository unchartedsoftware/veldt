package api

import (
	"errors"
	"fmt"
	"strconv"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/server/elastic"
)

type tileRes struct {
	Bins []byte `json:"bins"`
	Status string `json:"status"`
}

func parseTileCoord( params map[string]string ) ( *binning.TileCoord, error ) {
	x, ex := strconv.ParseUint( params["x"], 10, 32 )
	y, ey := strconv.ParseUint( params["y"], 10, 32 )
	z, ez := strconv.ParseUint( params["z"], 10, 32 )
	if ex == nil || ey == nil || ez == nil {
		return &binning.TileCoord{
			X: uint32( x ),
			Y: uint32( y ),
			Z: uint32( z ),
		}, nil
	}
	return nil, errors.New( "Unable to parse tile coordinate from URL" )
}

func handleTileErr( w http.ResponseWriter ) {
	// send error
	w.WriteHeader( 500 )
	fmt.Fprint( w, `{"status": "error"}` )
}

func tileHandler( c web.C, w http.ResponseWriter, r *http.Request ) {
	// set content type response header
	w.Header().Set( "Content-Type", "application/octet-stream" )
	// parse tile coord from URL
	tile, parseErr := parseTileCoord( c.URLParams )
	if parseErr != nil {
		handleTileErr( w )
		return
	}
	// extract tile data
	bins, tileErr := elastic.GetTile( tile )
	if tileErr != nil {
		handleTileErr( w )
		return
	}
	// send response
	w.WriteHeader( 200 )
	w.Write( bins )
}
