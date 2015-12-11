package meta

import (
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/log"
)

// GetMetaByType returns a meta data response based on the provided hash and
// request object.
func GetMetaByType(metaHash string, metaReq *MetaRequest) *MetaResponse {
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
