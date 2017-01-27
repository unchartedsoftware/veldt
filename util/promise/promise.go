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
	mutex    sync.Mutex
}

// NewPromise instantiates and returns a new promise.
func NewPromise() *Promise {
	return &Promise{
		Chan:     make(chan error),
		count:    0,
		resolved: false,
		response: nil,
		mutex:    sync.Mutex{},
	}
}

// Wait returns a channel that the response will be passed once the promise is
// resolved.
func (p *Promise) Wait() error {
	p.mutex.Lock()
	if p.resolved {
		p.mutex.Unlock()
		runtime.Gosched()
		return p.response
	}
	p.count++
	p.mutex.Unlock()
	runtime.Gosched()
	return <-p.Chan
}

// Resolve waits the response and sends it to all clients waiting on the channel.
func (p *Promise) Resolve(res error) {
	p.mutex.Lock()
	defer runtime.Gosched()
	defer p.mutex.Unlock()
	if p.resolved {
		return
	}
	p.resolved = true
	p.response = res
	for i := 0; i < p.count; i++ {
		p.Chan <- res
	}
}
