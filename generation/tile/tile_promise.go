package tile

import (
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/promise"
)

var (
	promises = promise.NewMap()
)

func getTilePromise(tileHash string, tileReq *Request, storeReq *store.Request, tileGen Generator) chan error {
	p, ok := promises.Get(tileHash)
	if ok {
		return p.Wait()
	}
	p = promise.NewPromise()
	promises.Set(tileHash, p)
	go func() {
		err := generateAndStoreTile(tileHash, tileReq, storeReq, tileGen)
		p.Resolve(err)
		promises.Remove(tileHash)
	}()
	return p.Wait()
}
