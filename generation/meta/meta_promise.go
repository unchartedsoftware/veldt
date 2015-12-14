package meta

import (
	"runtime"
	"sync"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/util/log"
)

var (
	mutex        = sync.Mutex{}
	metaPromises = make(map[string]*promise.Promise)
)

func getSuccessPromise(metaReq *Request, meta []byte) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(getSuccessResponse(metaReq, meta))
	log.Infof("Resolved meta promise with len of bytes %d", len(meta))
	return p
}

func getMetaPromise(metaHash string, metaReq *Request) *promise.Promise {
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
		meta := GenerateAndStoreMeta(metaHash, metaReq)
		p.Resolve(meta)
		mutex.Lock()
		delete(metaPromises, metaHash)
		mutex.Unlock()
		runtime.Gosched()
	}()
	return p
}
