package promise

import (
	"runtime"
	"sync"
)

// Promise represents a channel that will be shared by a variable number of
// users.
type Promise struct {
	Chan     chan error
	count    int
	resolved bool
	response error
	mutex    sync.RWMutex
}

// NewPromise instantiates and returns a new promise.
func NewPromise() *Promise {
	return &Promise{
		Chan:     make(chan error),
		count:    0,
		resolved: false,
		response: nil,
		mutex:    sync.RWMutex{},
	}
}

// Wait returns a channel that the response will be passed once the promise is
// resolved.
func (p *Promise) Wait() chan error {
	p.mutex.RLock()
	if p.resolved {
		p.mutex.RUnlock()
		runtime.Gosched()
		go func() {
			p.Chan <- p.response
		}()
		return p.Chan
	}
	p.count++
	p.mutex.RUnlock()
	runtime.Gosched()
	return p.Chan
}

// Resolve waits the reponse and sends it to all clients waiting on the channel.
func (p *Promise) Resolve(res error) {
	p.mutex.Lock()
	if p.resolved {
		p.mutex.Unlock()
		runtime.Gosched()
		return
	}
	p.resolved = true
	p.response = res
	p.mutex.Unlock()
	runtime.Gosched()
	for i := 0; i < p.count; i++ {
		p.Chan <- res
	}
}
