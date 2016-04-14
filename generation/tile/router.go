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
	Coord  *binning.TileCoord     `json:"coord"`
	Type   string                 `json:"type"`
	Index  string                 `json:"index"`
	Store  string                 `json:"store"`
	Params map[string]interface{} `json:"params"`
}

// String returns the request formatted as a string.
func (r *Request) String() string {
	return fmt.Sprintf("%s/%s/%d/%d/%d",
		r.Type,
		r.Index,
		r.Coord.Z,
		r.Coord.X,
		r.Coord.Y)
}

func getTileHash(req *Request, gen Generator) string {
	tileParams := gen.GetParams()
	// create hashes array
	var hashes []string
	// add tile req hash first
	hash := fmt.Sprintf("tile:%s:%s:%s:%d:%d:%d",
		gen.GetHash(),
		req.Type,
		req.Index,
		req.Coord.Z,
		req.Coord.X,
		req.Coord.Y)
	hashes = append(hashes, hash)
	// add individual param hashes
	for _, p := range tileParams {
		// check if the value held by the typed interface is null (a typed interface itself is never null)
		if reflect.ValueOf(p).IsNil() {
			hashes = append(hashes, "-")
		} else {
			hashes = append(hashes, p.GetHash())
		}
	}
	return strings.Join(hashes, ":")
}

// generateAndStoreTile generates the tile and puts it in the store.
func generateAndStoreTile(hash string, req *Request, gen Generator) error {
	// queue the tile to be generated
	tile, err := queue(gen)
	if err != nil {
		return err
	}
	// add tile to store
	return store.Set(req.Store, hash, tile[0:])
}

// GenerateTile issues a generation request and returns an error when it has
// completed.
func GenerateTile(req *Request) error {
	// get parameters
	gen, err := GetGenerator(req)
	if err != nil {
		return err
	}
	// get tile hash
	hash := getTileHash(req, gen)
	// check if tile already exists in store
	exists, err := store.Exists(req.Store, hash)
	if err != nil {
		return err
	}
	// if it exists, return as success
	if exists {
		return nil
	}
	// otherwise, initiate the tiling job and return error
	return getTilePromise(hash, req, gen)
}

// GetTileFromStore returns a serialized tile from store.
func GetTileFromStore(req *Request) ([]byte, error) {
	// get parameters
	gen, err := GetGenerator(req)
	if err != nil {
		return nil, err
	}
	// get tile hash
	hash := getTileHash(req, gen)
	// get tile data from store
	return store.Get(req.Store, hash)
}
