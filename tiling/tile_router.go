package tiling

import (
	"github.com/unchartedsoftware/prism/store"
)

// GetTileByType returns a tile response based on the provided hash and request
// object.
func GetTileByType(tileHash string, tileReq *TileRequest) *TileResponse {
	// get tile generator by id
	gen, err := GetTileGeneratorByType(tileReq.Type)
	if err != nil {
		return getFailureResponse(tileReq, err)
	}
	// generate tile
	tileData, err := gen(tileReq)
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
