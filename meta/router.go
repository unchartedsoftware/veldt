package meta

import (
	"fmt"

	"github.com/unchartedsoftware/prism/store"
)

func generateAndStoreMeta(hash string, req *Request, gen Generator) error {
	// generate meta
	meta, err := gen.GetMeta(req.URI)
	if err != nil {
		return err
	}
	// add meta data to store
	return store.Set(req.Store, hash, meta[0:])
}

// GenerateMeta issues a generation request and returns an error when it has
// completed.
func GenerateMeta(req *Request) error {
	// get meta generator by id
	gen, err := GetGenerator(req)
	if err != nil {
		return err
	}
	hash := req.GetHash()
	// check if meta data already exists in store
	exists, err := store.Exists(req.Store, hash)
	if err != nil {
		return err
	}
	// if it exists, return as success
	if exists {
		return nil
	}
	// otherwise, generate the metadata and return error
	return getMetaPromise(hash, req, gen)
}

// GetMetaFromStore returns serialized meta data from store.
func GetMetaFromStore(req *Request) ([]byte, error) {
	// get meta generator by id
	gen, err := GetGenerator(req)
	if err != nil {
		return nil, err
	}
	hash := req.GetHash()
	// get meta data from store
	return store.Get(req.Store, hash)
}
