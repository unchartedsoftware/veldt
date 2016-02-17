package meta

import (
	"fmt"

	"github.com/unchartedsoftware/prism/store"
)

// Request represents a meta data request.
type Request struct {
	Type  string `json:"type"`
	Index string `json:"index"`
	Store string `json:"store"`
}

// String returns the request formatted as a string.
func (r *Request) String() string {
	return fmt.Sprintf("%s/%s",
		r.Type,
		r.Index)
}

// getMetaHash returns a unique hash string for the given type.
func getMetaHash(metaReq *Request, metaGen Generator) string {
	return fmt.Sprintf("meta:%s:%s:%s",
		metaGen.GetHash(),
		metaReq.Type,
		metaReq.Index)
}

// generateAndStoreMeta generates the meta data and puts it in the store.
func generateAndStoreMeta(metaHash string, metaReq *Request, metaGen Generator) error {
	// generate meta
	meta, err := metaGen.GetMeta()
	if err != nil {
		return err
	}
	// get store connection
	conn, err := store.GetConnection(metaReq.Store)
	if err != nil {
		return err
	}
	// add tile to store
	err = conn.Set(metaHash, meta)
	conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// GenerateMeta issues a generation request and returns an error when it has
// completed.
func GenerateMeta(metaReq *Request) error {
	// get meta generator by id
	metaGen, err := GetGenerator(metaReq)
	if err != nil {
		return err
	}
	metaHash := getMetaHash(metaReq, metaGen)
	// get store connection
	conn, err := store.GetConnection(metaReq.Store)
	if err != nil {
		return err
	}
	// check if meta exists in store
	exists, err := conn.Exists(metaHash)
	if err != nil {
		return err
	}
	// if it exists, return as success
	if exists {
		return nil
	}
	// otherwise, generate the metadata and return error
	return getMetaPromise(metaHash, metaReq, metaGen)
}

// GetMetaFromStore returns serialized meta data from store.
func GetMetaFromStore(metaReq *Request) ([]byte, error) {
	// get meta generator by id
	metaGen, err := GetGenerator(metaReq)
	if err != nil {
		return nil, err
	}
	metaHash := getMetaHash(metaReq, metaGen)
	// get store connection
	conn, err := store.GetConnection(metaReq.Store)
	if err != nil {
		return nil, err
	}
	// get meta data from store
	meta, err := conn.Get(metaHash)
	conn.Close()
	if err != nil {
		return nil, err
	}
	return meta, nil
}
