package meta

import (
	"fmt"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/log"
)

// Request represents a meta data request.
type Request struct {
	Type     string `json:"type"`
	Index    string `json:"index"`
	Endpoint string `json:"endpoint"`
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
func GetMeta(metaReq *Request) *promise.Promise {
	metaHash := getMetaHash(metaReq)
	// check if meta exists in store
	exists, err := store.Exists(metaHash)
	if err != nil {
		log.Warn(err)
	}
	// if it exists, return success promise
	if exists {
		meta, err := store.Get(metaHash)
		if err == nil && meta != nil {
			// if no error, return promise
			return getSuccessPromise(metaReq, meta)
		}
	}
	// otherwise, generate the metadata and return promise
	return getMetaPromise(metaHash, metaReq)
}

// GenerateAndStoreMeta returns a meta data response based on the provided hash and
// request object.
func GenerateAndStoreMeta(metaHash string, metaReq *Request) *Response {
	// get meta generator by id
	gen, err := GetGeneratorByType(metaReq.Type)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	// generate meta
	meta, err := gen(metaReq)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	// add tile to store
	err = store.Set(metaHash, meta)
	if err != nil {
		log.Warn(err)
	}
	return getSuccessResponse(metaReq, meta)
}
