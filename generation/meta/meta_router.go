package meta

import (
	"fmt"

	"github.com/unchartedsoftware/prism/store"
)

// Request represents a meta data request.
type Request struct {
	Type     string `json:"type"`
	Index    string `json:"index"`
	Endpoint string `json:"endpoint"`
}

// String returns the request formatted as a string.
func (r *Request) String() string {
	return fmt.Sprintf("%s/%s/%s",
		r.Endpoint,
		r.Index,
		r.Type)
}

// getMetaHash returns a unique hash string for the given type.
func getMetaHash(metaReq *Request) string {
	return fmt.Sprintf("%s:%s:%s:meta",
		metaReq.Endpoint,
		metaReq.Index,
		metaReq.Type)
}

// generateAndStoreMeta generates the meta data and puts it in the store.
func generateAndStoreMeta(metaHash string, metaReq *Request, storeReq *store.Request) error {
	// get meta generator by id
	metaGen, err := GetGenerator(metaReq)
	if err != nil {
		return err
	}
	// generate meta
	meta, err := metaGen.GetMeta(metaReq)
	if err != nil {
		return err
	}
	// get store connection
	conn, err := store.GetConnection(storeReq)
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

// GenerateMeta returns a promise which will be fulfilled when the meta
// generation has completed.
func GenerateMeta(metaReq *Request, storeReq *store.Request) error {
	metaHash := getMetaHash(metaReq)
	// get store connection
	conn, err := store.GetConnection(storeReq)
	if err != nil {
		return err
	}
	// check if meta exists in store
	exists, err := conn.Exists(metaHash)
	if err != nil {
		return err
	}
	// if it exists, return success promise
	if exists {
		return nil
	}
	// otherwise, generate the metadata and return promise
	return getMetaPromise(metaHash, metaReq, storeReq)
}

// GetMetaFromStore returns serialized meta data from store.
func GetMetaFromStore(metaReq *Request, storeReq *store.Request) ([]byte, error) {
	// get meta hash
	metaHash := getMetaHash(metaReq)
	// get store connection
	conn, err := store.GetConnection(storeReq)
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
