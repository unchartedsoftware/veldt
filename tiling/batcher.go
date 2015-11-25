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
}

// TileResponse represents the tile response data
type TileResponse struct {
	TileCoord binning.TileCoord `json:"tilecoord"`
	Type      string            `json:"type"`
	Success   bool              `json:"success"`
}

func getFailureResponse(tileReq *TileRequest) *TileResponse {
	return &TileResponse{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Success:   false,
	}
}

func getSuccessResponse(tileReq *TileRequest) *TileResponse {
	return &TileResponse{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Success:   true,
	}
}

// GetTileHash returns a unique hash string for the given tile and type.
func GetTileHash(tileReq *TileRequest) string {
	return fmt.Sprintf("%s-%d-%d-%d",
		tileReq.Type,
		tileReq.TileCoord.X,
		tileReq.TileCoord.Y,
		tileReq.TileCoord.Z)
}

// GetTile will return a tile response channel that can be used to determine when a tile is ready.
func GetTile(tileReq *TileRequest) (*promise.Promise, error) {
	// get hash for tile
	tileHash := GetTileHash(tileReq)
	// check if tile exists in store
	exists, err := store.Exists(tileHash)
	if err != nil {
		log.Warn(err)
	}
	// if it exists, return success promise
	if exists {
		return GetSuccessPromise(tileReq), nil
	}
	// otherwise, initiate the tiling job and return promise
	return GetTilePromise(tileHash, tileReq), nil
}
