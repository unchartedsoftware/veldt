package meta

import (
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/promise"
)

var (
	promises = promise.NewMap()
)

func getMetaPromise(metaHash string, metaReq *Request, storeReq *store.Request) error {
	p, ok := promises.GetOrCreate(metaHash)
	if ok {
		// promise already existed, return it
		return p.Wait()
	}
	// promise had to be created, generate meta data
	go func() {
		err := generateAndStoreMeta(metaHash, metaReq, storeReq)
		p.Resolve(err)
		promises.Remove(metaHash)
	}()
	return p.Wait()
}
