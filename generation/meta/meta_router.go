package meta

import (
	"fmt"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/log"
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

// Response represents the meta data response.
type Response struct {
	Meta     []byte `json:"meta"`
	Type     string `json:"type"`
	Index    string `json:"index"`
	Endpoint string `json:"endpoint"`
	Success  bool   `json:"success"`
	Error    error  `json:"-"` // ignore field
}

func getFailureResponse(metaReq *Request, err error) *Response {
	return &Response{
		Type:     metaReq.Type,
		Index:    metaReq.Index,
		Endpoint: metaReq.Endpoint,
		Success:  false,
		Error:    err,
	}
}

func getSuccessResponse(metaReq *Request, meta []byte) *Response {
	return &Response{
		Meta:     meta,
		Type:     metaReq.Type,
		Index:    metaReq.Index,
		Endpoint: metaReq.Endpoint,
		Success:  true,
		Error:    nil,
	}
}

// getMetaHash returns a unique hash string for the given type.
func getMetaHash(metaReq *Request) string {
	return fmt.Sprintf("%s:%s:%s:meta",
		metaReq.Endpoint,
		metaReq.Index,
		metaReq.Type)
}

// GetMeta returns a promise which will be fulfilled when the meta generation
// has completed.
func GetMeta(metaReq *Request, storeReq *store.Request) *promise.Promise {
	metaHash := getMetaHash(metaReq)
	// get store connection
	conn, err := store.GetConnection(storeReq)
	if err != nil {
		// only log warning, we can still generate meta, we just can't store it
		log.Warn(err)
		return getMetaPromise(metaHash, metaReq, storeReq)
	}
	// check if meta exists in store
	exists, err := conn.Exists(metaHash)
	if err != nil {
		// only log warning, we can still generate meta, we just can't store it
		log.Warn(err)
		conn.Close()
		return getMetaPromise(metaHash, metaReq, storeReq)
	}
	// if it exists, return success promise
	if exists {
		meta, err := conn.Get(metaHash)
		conn.Close()
		if err == nil && meta != nil {
			// if no error, return promise
			return getSuccessPromise(metaReq, meta)
		}
	}
	// otherwise, generate the metadata and return promise
	return getMetaPromise(metaHash, metaReq, storeReq)
}

// GenerateAndStoreMeta returns a meta data response based on the provided hash and
// request object.
func GenerateAndStoreMeta(metaHash string, metaReq *Request, storeReq *store.Request) *Response {
	// get meta generator by id
	metaGen, err := GetGenerator(metaReq)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	// generate meta
	meta, err := metaGen.GetMeta(metaReq)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	// get store connection
	conn, err := store.GetConnection(storeReq)
	if err != nil {
		// only log warning, we can still generate meta, we just can't store it
		log.Warn(err)
	} else {
		// add tile to store
		err := conn.Set(metaHash, meta)
		conn.Close()
		if err != nil {
			log.Warn(err)
		}
	}
	return getSuccessResponse(metaReq, meta)
}
