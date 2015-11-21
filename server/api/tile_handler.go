package api

import (
	"errors"
	"fmt"
	"strconv"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tiling/elastic"
)

type tileRes struct {
	Bins []byte `json:"bins"`
	Status string `json:"status"`
}

// TileRequest represents the tile type and tile coord
type TileRequest struct {
	TileCoord binning.TileCoord
	Type string
}

func parseTileParams( params map[string]string ) ( *TileRequest, error ) {
	tileType := params["type"]
	x, ex := strconv.ParseUint( params["x"], 10, 32 )
	y, ey := strconv.ParseUint( params["y"], 10, 32 )
	z, ez := strconv.ParseUint( params["z"], 10, 32 )
	if ex == nil || ey == nil || ez == nil {
		return &TileRequest {
			TileCoord: binning.TileCoord{
				X: uint32( x ),
				Y: uint32( y ),
				Z: uint32( z ),
			},
			Type: tileType,
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
	w.Header().Set( "Content-Type", "application/json" )
	// parse tile coord from URL
	tileReq, parseErr := parseTileParams( c.URLParams )
	if parseErr != nil {
		handleTileErr( w )
		return
	}
	var tileData []byte
	var tileErr error
	// get tile based on type
	if tileReq.Type == "topiccount" {
		tileData, tileErr = elastic.GetTopicCountTile( &tileReq.TileCoord )
	} else {
		tileData, tileErr = elastic.GetHeatmapTile( &tileReq.TileCoord )
	}
	if tileErr != nil {
		handleTileErr( w )
		return
	}
	// send response
	w.WriteHeader( 200 )
	w.Write( tileData )
}
