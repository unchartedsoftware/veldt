package meta

import (
	"runtime"
	"sync"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/store"
)

var (
	mutex        = sync.Mutex{}
	metaPromises = make(map[string]*promise.Promise)
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

func getMetaPromise(metaHash string, metaReq *Request, storeReq *store.Request) *promise.Promise {
	mutex.Lock()
	p, ok := metaPromises[metaHash]
	if ok {
		mutex.Unlock()
		runtime.Gosched()
		return p
	}
	p = promise.NewPromise()
	metaPromises[metaHash] = p
	mutex.Unlock()
	runtime.Gosched()
	go func() {
		err := generateAndStoreMeta(metaHash, metaReq, storeReq)
		p.Resolve(err)
		mutex.Lock()
		delete(metaPromises, metaHash)
		mutex.Unlock()
		runtime.Gosched()
	}()
	return p
}
