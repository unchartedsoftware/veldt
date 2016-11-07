package tile

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/unchartedsoftware/prism/store"
)

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
	gen, err := GetGenerator(req.Type)
	if err != nil {
		return err
	}
	// get tile hash
	hash := req.GetHash()
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
	hash := req.GetHash()
	// get tile data from store
	return store.Get(req.Store, hash)
}
