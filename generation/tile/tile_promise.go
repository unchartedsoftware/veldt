package tile

import (
	"runtime"
	"sync"

	"github.com/fanliao/go-promise"
)

var (
	mutex        = sync.Mutex{}
	tilePromises = make(map[string]*promise.Promise)
)

func getSuccessPromise(tileReq *Request) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(getSuccessResponse(tileReq))
	return p
}

func getFailurePromise(tileReq *Request, err error) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(getFailureResponse(tileReq, err))
	return p
}

func getTilePromise(tileHash string, tileReq *Request) *promise.Promise {
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
		tileRes := GenerateAndStoreTile(tileHash, tileReq)
		p.Resolve(tileRes)
		mutex.Lock()
		delete(tilePromises, tileHash)
		mutex.Unlock()
		runtime.Gosched()
	}()
	return p
}
