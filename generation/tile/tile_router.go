package tile

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/store"
)

// Request represents the tile type and tile coord.
type Request struct {
	Coord    *binning.TileCoord     `json:"coord"`
	Type     string                 `json:"type"`
	Index    string                 `json:"index"`
	Endpoint string                 `json:"endpoint"`
	Params   map[string]interface{} `json:"params"`
}

// String returns the request formatted as a string.
func (r *Request) String() string {
	return fmt.Sprintf("%s/%s/%s/%d/%d/%d",
		r.Endpoint,
		r.Index,
		r.Type,
		r.Coord.Z,
		r.Coord.X,
		r.Coord.Y)
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
		tileReq.Coord.X,
		tileReq.Coord.Y,
		tileReq.Coord.Z)
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

// generateAndStoreTile generates the tile and puts it in the store.
func generateAndStoreTile(tileHash string, tileReq *Request, storeReq *store.Request, tileGen Generator) error {
	// generate the tile
	tileData, err := tileGen.GetTile(tileReq)
	if err != nil {
		return err
	}
	// get store connection
	conn, err := store.GetConnection(storeReq)
	if err != nil {
		return err
	}
	// add tile to store
	err = conn.Set(tileHash, tileData[0:])
	conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// GenerateTile issues a generation request and returns an error when it has
// completed.
func GenerateTile(tileReq *Request, storeReq *store.Request) error {
	// get parameters
	tileGen, err := GetGenerator(tileReq)
	if err != nil {
		return err
	}
	// get store connection
	conn, err := store.GetConnection(storeReq)
	if err != nil {
		return err
	}
	// get tile hash
	tileHash := getTileHash(tileReq, tileGen)
	// check if tile exists in store
	exists, err := conn.Exists(tileHash)
	conn.Close()
	if err != nil {
		return err
	}
	// if it exists, return as success
	if exists {
		return nil
	}
	// otherwise, initiate the tiling job and return error
	return getTilePromise(tileHash, tileReq, storeReq, tileGen)
}

// GetTileFromStore returns a serialized tile from store.
func GetTileFromStore(tileReq *Request, storeReq *store.Request) ([]byte, error) {
	// get parameters
	tileGen, err := GetGenerator(tileReq)
	if err != nil {
		return nil, err
	}
	// get tile hash
	tileHash := getTileHash(tileReq, tileGen)
	// get store connection
	conn, err := store.GetConnection(storeReq)
	if err != nil {
		return nil, err
	}
	// get tile data from store
	tile, err := conn.Get(tileHash)
	conn.Close()
	if err != nil {
		return nil, err
	}
	return tile, nil
}
