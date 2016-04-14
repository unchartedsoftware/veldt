package tile

import (
	"github.com/unchartedsoftware/prism/util/promise"
)

var (
	promises = promise.NewMap()
)

func getTilePromise(hash string, req *Request, gen Generator) error {
	p, exists := promises.GetOrCreate(hash)
	if exists {
		// promise already existed, return it
		return p.Wait()
	}
	// promise had to be created, generate tile
	go func() {
		err := generateAndStoreTile(hash, req, gen)
		p.Resolve(err)
		promises.Remove(hash)
	}()
	return p.Wait()
}
