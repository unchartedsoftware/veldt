package meta

import (
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/promise"
)

var (
	promises = promise.NewMap()
)

func getMetaPromise(metaHash string, metaReq *Request, storeReq *store.Request) chan error {
	p, ok := promises.Get(metaHash)
	if ok {
		return p.Wait()
	}
	p = promise.NewPromise()
	promises.Set(metaHash, p)
	go func() {
		err := generateAndStoreMeta(metaHash, metaReq, storeReq)
		p.Resolve(err)
		promises.Remove(metaHash)
	}()
	return p.Wait()
}
