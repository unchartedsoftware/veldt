package tiling

import (
	"errors"

	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/tiling/elastic"
)

func generateTileByType(tileReq *TileRequest) ([]byte, error) {
	switch tileReq.Type {
	case "topiccount":
		return elastic.GetTopicCountTile(tileReq.Endpoint, tileReq.Index, &tileReq.TileCoord)
	case "heatmap":
		return elastic.GetHeatmapTile(tileReq.Endpoint, tileReq.Index, &tileReq.TileCoord)
	default:
		return nil, errors.New("Tiling type not recognized")
	}
}

// GetTileByType returns a tile response based on the provided hash and request object.
func GetTileByType(tileHash string, tileReq *TileRequest) *TileResponse {
	// generate tile by type
	tileData, err := generateTileByType(tileReq)
	if err != nil {
		return getFailureResponse(tileReq, err)
	}
	// add tile to store
	err = store.Set(tileHash, tileData[0:])
	if err != nil {
		return getFailureResponse(tileReq, err)
	}
	return getSuccessResponse(tileReq)
}
