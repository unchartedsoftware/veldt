package tiling

import (
	"errors"
	"fmt"
    "sync"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/tiling/elastic"
)

// TileRequest represents the tile type and tile coord
type TileRequest struct {
	TileCoord binning.TileCoord `json:"tilecoord"`
	Type string `json:"type"`
}

// TileResponse represents the tile response data
type TileResponse struct {
	TileCoord binning.TileCoord `json:"tilecoord"`
	Type string `json:"type"`
	Success bool `json:"success"`
}

var mutex = sync.Mutex{}
var promiseMap = make( map[string]chan TileResponse )

func getFailureResponse( tileReq *TileRequest ) TileResponse {
	return TileResponse{
		TileCoord: tileReq.TileCoord,
		Type: tileReq.Type,
		Success: false,
	}
}

func getSuccessResponse( tileReq *TileRequest ) TileResponse {
	return TileResponse{
		TileCoord: tileReq.TileCoord,
		Type: tileReq.Type,
		Success: true,
	}
}

func generateTileByType( tileReq *TileRequest ) ( []byte, error ){
	switch tileReq.Type {
		case "topiccount":
			return elastic.GetTopicCountTile( &tileReq.TileCoord )
		case "heatmap":
			return elastic.GetHeatmapTile( &tileReq.TileCoord )
		default:
			return nil, errors.New("Tiling type not recognized")
	}
}

func generateTile( tileHash string, tileReq *TileRequest ) TileResponse {
	// generate tile by type
	tileData, err := generateTileByType( tileReq )
	if err != nil {
		return getFailureResponse( tileReq )
	}
	// add tile to store
	store.Set( tileHash, tileData )
	return getSuccessResponse( tileReq )
}

func getTileHash( tileReq *TileRequest ) string {
	return fmt.Sprintf( "%s-%d-%d-%d-",
		tileReq.Type,
		tileReq.TileCoord.X,
		tileReq.TileCoord.Y,
		tileReq.TileCoord.Z )
}

func getTilePromise( tileHash string, tileReq *TileRequest ) chan TileResponse {
    promise := make( chan TileResponse )
    go func () {
		promise <- generateTile( tileHash, tileReq )
	}()
    return promise
}

func getSuccessPromise( tileReq *TileRequest ) chan TileResponse {
	promise := make( chan TileResponse )
    go func () {
		promise <- getSuccessResponse( tileReq )
	}()
    return promise
}

func getFailurePromise( tileReq *TileRequest ) chan TileResponse {
	promise := make( chan TileResponse )
    go func () {
		promise <- getFailureResponse( tileReq )
	}()
    return promise
}

// GetTile will return a tile response channel that can be used to determine when a tile is ready.
func GetTile( tileReq *TileRequest ) chan TileResponse {
	// get hash for tile
	tileHash := getTileHash( tileReq )
	// check if tile exists in store
	exists, err := store.Exists( tileHash )
	if err != nil {
		// on error return failure promise
		return getFailurePromise( tileReq )
	}
	// if it exists, return success promise
	if exists {
		return getSuccessPromise( tileReq )
	}
	// otherwise, initiate the tiling job
	mutex.Lock()
	if promiseMap[ tileHash ] == nil {
		promiseMap[ tileHash ] = getTilePromise( tileHash, tileReq )
	}
	mutex.Unlock()
	return promiseMap[ tileHash ]
}
