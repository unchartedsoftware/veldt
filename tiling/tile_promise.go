package tiling

import (
	"sync"

	"github.com/fanliao/go-promise"

	"github.com/unchartedsoftware/prism/util/log"
)

var mutex = sync.Mutex{}
var tilePromises = make(map[string]*promise.Promise)

func deletePromise(tileHash string) {
	mutex.Lock()
	delete(tilePromises, tileHash)
	mutex.Unlock()
}

// GetTilePromise will return a promise for an existing tiling request, or initiate a tiling request and return a promise.
func GetTilePromise(tileHash string, tileReq *TileRequest) *promise.Promise {
	mutex.Lock()
	p, ok := tilePromises[tileHash]
	if ok {
		mutex.Unlock()
		return p
	}
	p = promise.NewPromise()
	tilePromises[tileHash] = p
	mutex.Unlock()
	go func() {
		tileRes, err := GetTileByType(tileHash, tileReq)
		if err != nil {
			log.Error(err)
		}
		p.Resolve(tileRes)
		deletePromise(tileHash)
	}()
	return p
}

// GetSuccessPromise returns a successful request promise that immediately resolves.
func GetSuccessPromise(tileReq *TileRequest) *promise.Promise {
	p := promise.NewPromise()
	p.Resolve(getSuccessResponse(tileReq))
	return p
}
