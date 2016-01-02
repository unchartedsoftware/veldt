package tile

import (
	"runtime"
	"sync"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/store"
)

var (
	mutex        = sync.Mutex{}
	tilePromises = make(map[string]*promise.Promise)
)

func getSuccessPromise() *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(nil)
	return p
}

func getFailurePromise(err error) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(err)
	return p
}

func getTilePromise(tileHash string, tileReq *Request, storeReq *store.Request, tileGen Generator) *promise.Promise {
	mutex.Lock()
	p, ok := tilePromises[tileHash]
	if ok {
		mutex.Unlock()
		runtime.Gosched()
		return p
	}
	p = promise.NewPromise()
	tilePromises[tileHash] = p
	mutex.Unlock()
	runtime.Gosched()
	go func() {
		err := generateAndStoreTile(tileHash, tileReq, storeReq, tileGen)
		p.Resolve(err)
		mutex.Lock()
		delete(tilePromises, tileHash)
		mutex.Unlock()
		runtime.Gosched()
	}()
	return p
}
