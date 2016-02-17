package meta

import (
	"github.com/unchartedsoftware/prism/util/promise"
)

var (
	promises = promise.NewMap()
)

func getMetaPromise(metaHash string, metaReq *Request, metaGen Generator) error {
	p, exists := promises.GetOrCreate(metaHash)
	if exists {
		// promise already existed, return it
		return p.Wait()
	}
	// promise had to be created, generate meta data
	go func() {
		err := generateAndStoreMeta(metaHash, metaReq, metaGen)
		p.Resolve(err)
		promises.Remove(metaHash)
	}()
	return p.Wait()
}
