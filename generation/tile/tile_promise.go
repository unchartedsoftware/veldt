package tile

import (
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/util/promise"
)

var (
	promises = promise.NewMap()
)

func getTilePromise(tileHash string, tileReq *Request, storeReq *store.Request, tileGen Generator) error {
	p, ok := promises.GetOrCreate(tileHash)
	if ok {
		// promise already existed, return it
		return p.Wait()
	}
	// promise had to be created, generate tile
	go func() {
		err := generateAndStoreTile(tileHash, tileReq, storeReq, tileGen)
		p.Resolve(err)
		promises.Remove(tileHash)
	}()
	return p.Wait()
}
