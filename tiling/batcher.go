package tiling

import (
	"fmt"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/log"
)

// TileRequest represents the tile type and tile coord
type TileRequest struct {
	TileCoord binning.TileCoord `json:"tilecoord"`
	Type      string            `json:"type"`
	Index     string            `json:"index"`
	Endpoint  string            `json:"endpoint"`
}

// TileResponse represents the tile response data
type TileResponse struct {
	TileCoord binning.TileCoord `json:"tilecoord"`
	Type      string            `json:"type"`
	Index     string            `json:"index"`
	Endpoint  string            `json:"endpoint"`
	Success   bool              `json:"success"`
	Error     error             `json:"-"` // ignore field
}

func getFailureResponse(tileReq *TileRequest, err error) *TileResponse {
	return &TileResponse{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Index:     tileReq.Index,
		Endpoint:  tileReq.Endpoint,
		Success:   false,
		Error:     err,
	}
}

func getSuccessResponse(tileReq *TileRequest) *TileResponse {
	return &TileResponse{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Index:     tileReq.Index,
		Endpoint:  tileReq.Endpoint,
		Success:   true,
		Error:     nil,
	}
}

// GetTileHash returns a unique hash string for the given tile and type.
func GetTileHash(tileReq *TileRequest) string {
	return fmt.Sprintf("%s:%s:%s:%d-%d-%d",
		tileReq.Endpoint,
		tileReq.Index,
		tileReq.Type,
		tileReq.TileCoord.X,
		tileReq.TileCoord.Y,
		tileReq.TileCoord.Z)
}

// GetTile will return a tile response channel that can be used to determine when a tile is ready.
func GetTile(tileReq *TileRequest) *promise.Promise {
	// get hash for tile
	tileHash := GetTileHash(tileReq)
	// check if tile exists in store
	exists, err := store.Exists(tileHash)
	if err != nil {
		log.Warn(err)
	}
	// if it exists, return success promise
	if exists {
		return GetSuccessPromise(tileReq)
	}
	// otherwise, initiate the tiling job and return promise
	return GetTilePromise(tileHash, tileReq)
}
