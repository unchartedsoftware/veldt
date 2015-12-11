package meta

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/log"
)

var (
	mutex        = sync.Mutex{}
	metaPromises = make(map[string]*promise.Promise)
)

// MetaRequest represents a meta data request.
type MetaRequest struct {
	Type     string `json:"type"`
	Index    string `json:"index"`
	Endpoint string `json:"endpoint"`
}

// MetaResponse represents the meta data response.
type MetaResponse struct {
	Meta     []byte `json:"meta"`
	Type     string `json:"type"`
	Index    string `json:"index"`
	Endpoint string `json:"endpoint"`
	Success  bool   `json:"success"`
	Error    error  `json:"-"` // ignore field
}

func getFailureResponse(metaReq *MetaRequest, err error) *MetaResponse {
	return &MetaResponse{
		Type:     metaReq.Type,
		Index:    metaReq.Index,
		Endpoint: metaReq.Endpoint,
		Success:  false,
		Error:    err,
	}
}

func getSuccessResponse(metaReq *MetaRequest, meta []byte) *MetaResponse {
	return &MetaResponse{
		Meta:     meta,
		Type:     metaReq.Type,
		Index:    metaReq.Index,
		Endpoint: metaReq.Endpoint,
		Success:  true,
		Error:    nil,
	}
}

func getSuccessPromise(metaReq *MetaRequest, meta []byte) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(getSuccessResponse(metaReq, meta))
	return p
}

func getMetaPromise(metaHash string, metaReq *MetaRequest) *promise.Promise {
	mutex.Lock()
	p, ok := metaPromises[metaHash]
	if ok {
		mutex.Unlock()
		runtime.Gosched()
		return p
	}
	p = promise.NewPromise()
	metaPromises[metaHash] = p
	mutex.Unlock()
	runtime.Gosched()
	go func() {
		meta := GetMetaByType(metaHash, metaReq)
		p.Resolve(meta)
		mutex.Lock()
		delete(metaPromises, metaHash)
		mutex.Unlock()
		runtime.Gosched()
	}()
	return p
}

// GetMetaHash returns a unique hash string for the given type.
func GetMetaHash(metaReq *MetaRequest) string {
	return fmt.Sprintf("%s:%s:%s:meta",
		metaReq.Endpoint,
		metaReq.Index,
		metaReq.Type)
}

// GetMeta will return a promise for existing meta data, an existing meta
// request, or a new meta request.
func GetMeta(metaReq *MetaRequest) *promise.Promise {
	metaHash := GetMetaHash(metaReq)
	// check if meta exists in store
	exists, err := store.Exists(metaHash)
	if err != nil {
		log.Warn(err)
	}
	// if it exists, return success promise
	if exists {
		meta, err := store.Get(metaHash)
		if err != nil {
			// if err, generate it on the fly
			return getMetaPromise(metaHash, metaReq)
		}
		return getSuccessPromise(metaReq, meta)
	}
	// otherwise, generate the metadata and return promise
	return getMetaPromise(metaHash, metaReq)
}
