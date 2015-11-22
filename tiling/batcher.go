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
	Type      string            `json:"type"`
}

// TileResponse represents the tile response data
type TileResponse struct {
	TileCoord binning.TileCoord `json:"tilecoord"`
	Type      string            `json:"type"`
	Success   bool              `json:"success"`
	Error     error             `json:"-"`
}

var mutex = sync.Mutex{}

func getFailureResponse(tileReq *TileRequest, err error) *TileResponse {
	return &TileResponse{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Success:   false,
		Error:     err,
	}
}

func getSuccessResponse(tileReq *TileRequest) *TileResponse {
	return &TileResponse{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Success:   true,
		Error:     nil,
	}
}

func generateTileByType(tileReq *TileRequest) ([]byte, error) {
	switch tileReq.Type {
	case "topiccount":
		return elastic.GetTopicCountTile(&tileReq.TileCoord)
	case "heatmap":
		return elastic.GetHeatmapTile(&tileReq.TileCoord)
	default:
		return nil, errors.New("Tiling type not recognized")
	}
}

func generateTile(tileHash string, tileReq *TileRequest) *TileResponse {
	// generate tile by type
	tileData, err := generateTileByType(tileReq)
	if err != nil {
		return getFailureResponse(tileReq, err)
	}
	// add tile to store
	store.Set(tileHash, tileData)
	return getSuccessResponse(tileReq)
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
func GetTile(tileReq *TileRequest) *TileResponse {
	// get hash for tile
	tileHash := GetTileHash(tileReq)
	// check if tile exists in store
	exists, err := store.Exists(tileHash)
	if err != nil {
		fmt.Println(err)
	}
	// if it exists, return success promise
	if exists {
		return getSuccessResponse(tileReq)
	}
	// otherwise, initiate the tiling job
	return generateTile(tileHash, tileReq)
}
