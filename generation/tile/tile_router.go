package tile

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/log"
)

// Request represents the tile type and tile coord.
type Request struct {
	TileCoord *binning.TileCoord     `json:"tilecoord"`
	Type      string                 `json:"type"`
	Index     string                 `json:"index"`
	Endpoint  string                 `json:"endpoint"`
	Params    map[string]interface{} `json:"params"`
}

// Response represents the tile response data.
type Response struct {
	TileCoord *binning.TileCoord     `json:"tilecoord"`
	Type      string                 `json:"type"`
	Index     string                 `json:"index"`
	Endpoint  string                 `json:"endpoint"`
	Params    map[string]interface{} `json:"params"`
	Success   bool                   `json:"success"`
	Error     error                  `json:"-"` // ignore field
}

func getFailureResponse(tileReq *Request, err error) *Response {
	return &Response{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Index:     tileReq.Index,
		Endpoint:  tileReq.Endpoint,
		Params:    tileReq.Params,
		Success:   false,
		Error:     err,
	}
}

func getSuccessResponse(tileReq *Request) *Response {
	return &Response{
		TileCoord: tileReq.TileCoord,
		Type:      tileReq.Type,
		Index:     tileReq.Index,
		Endpoint:  tileReq.Endpoint,
		Params:    tileReq.Params,
		Success:   true,
		Error:     nil,
	}
}

func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}

func getTileHash(tileReq *Request, tileGen Generator) string {
	tileParams := tileGen.GetParams()
	// create hashes array
	var hashes []string
	// add tile req hash first
	hash := fmt.Sprintf("%s:%s:%s:%d:%d:%d",
		tileReq.Endpoint,
		tileReq.Index,
		tileReq.Type,
		tileReq.TileCoord.X,
		tileReq.TileCoord.Y,
		tileReq.TileCoord.Z)
	hashes = append(hashes, hash)
	// add individual param hashes
	for _, p := range tileParams {
		if isNil(p) {
			hashes = append(hashes, "-")
		} else {
			hashes = append(hashes, p.GetHash())
		}
	}
	return strings.Join(hashes, ":")
}

// GetTile returns a promise which will be fulfilled when the tile generation
// has completed and the tile is ready.
func GetTile(tileReq *Request) *promise.Promise {
	// get parameters
	tileGen, err := GetGenerator(tileReq)
	if err != nil {
		return getFailurePromise(tileReq, err)
	}
	// get tile hash
	tileHash := getTileHash(tileReq, tileGen)
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
	return getTilePromise(tileHash, tileReq, tileGen)
}

// GenerateAndStoreTile generates a new tile and then puts it in the store.
func GenerateAndStoreTile(tileHash string, tileReq *Request, tileGen Generator) *Response {
	// generate the tile
	tileData, err := tileGen.GetTile(tileReq)
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

// GetTileFromStore returns a serialized tile from store.
func GetTileFromStore(tileReq *Request) ([]byte, error) {
	// get parameters
	tileGen, err := GetGenerator(tileReq)
	if err != nil {
		return nil, err
	}
	// get tile hash
	tileHash := getTileHash(tileReq, tileGen)
	// get tile data from store
	tile, err := store.Get(tileHash)
	if tile == nil || err != nil {
		return nil, err
	}
	return tile, nil
}
