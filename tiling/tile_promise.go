package tiling

import (
	"runtime"
	"sync"

	"github.com/fanliao/go-promise"
)

var (
	mutex        = sync.Mutex{}
	tilePromises = make(map[string]*promise.Promise)
)

// GetTilePromise will return a promise for an existing tiling request, or initiate a tiling request and return a promise.
func GetTilePromise(tileHash string, tileReq *TileRequest) *promise.Promise {
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
		tileRes := GetTileByType(tileHash, tileReq)
		p.Resolve(tileRes)
		mutex.Lock()
		delete(tilePromises, tileHash)
		mutex.Unlock()
		runtime.Gosched()
	}()
	return p
}

// GetSuccessPromise returns a successful request promise that immediately resolves.
func GetSuccessPromise(tileReq *TileRequest) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(getSuccessResponse(tileReq))
	return p
}
