package meta

import (
	"github.com/unchartedsoftware/prism/util/promise"
)

var (
	promises = promise.NewMap()
)

func getMetaPromise(hash string, req *Request, gen Generator) error {
	p, exists := promises.GetOrCreate(hash)
	if exists {
		// promise already existed, return it
		return p.Wait()
	}
	// promise had to be created, generate meta data
	go func() {
		err := generateAndStoreMeta(hash, req, gen)
		p.Resolve(err)
		promises.Remove(hash)
	}()
	return p.Wait()
}
