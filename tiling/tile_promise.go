package tiling

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/log"
)

var (
	mutex        = sync.Mutex{}
	tilePromises = make(map[string]*promise.Promise)
)

// TileRequest represents the tile type and tile coord.
type TileRequest struct {
	TileCoord binning.TileCoord `json:"tilecoord"`
	Type      string            `json:"type"`
	Index     string            `json:"index"`
	Endpoint  string            `json:"endpoint"`
}

// TileResponse represents the tile response data.
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

func getSuccessPromise(tileReq *TileRequest) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(getSuccessResponse(tileReq))
	return p
}

func getTilePromise(tileHash string, tileReq *TileRequest) *promise.Promise {
	mutex.Lock()
	p, ok := tilePromises[tileHash]
	if ok {
		mutex.Unlock()
		runtime.Gosched()
		return p
	}
	p = promise.NewPromise()
	tilePromises[tileHash] = p
	mutex.Unlock()
	runtime.Gosched()
	go func() {
		tileRes := GetTileByType(tileHash, tileReq)
		p.Resolve(tileRes)
		mutex.Lock()
		delete(tilePromises, tileHash)
		mutex.Unlock()
		runtime.Gosched()
	}()
	return p
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

// GetTile will return a promise for existing tile data, an existing tiling
// request, or a new tiling request.
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
		return getSuccessPromise(tileReq)
	}
	// otherwise, initiate the tiling job and return promise
	return GetTilePromise(tileHash, tileReq)
}
