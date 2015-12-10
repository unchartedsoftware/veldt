package tiling

import (
	"encoding/json"
	
	"github.com/unchartedsoftware/prism/store"
)

// GetMetaByType returns a meta data response based on the provided hash and
// request object.
func GetMetaByType(metaHash string, metaReq *MetaRequest) *MetaResponse {
	// get meta generator by id
	gen, err := GetMetaGeneratorByType(metaReq.Type)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	// generate meta
	meta, err := gen(metaReq)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	// marshal data
	bytes, err := json.Marshal(meta)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	// add tile to store
	err = store.Set(metaHash, bytes)
	if err != nil {
		return getFailureResponse(metaReq, err)
	}
	return getSuccessResponse(metaReq, meta)
}
