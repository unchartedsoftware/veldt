package promise

import (
	"runtime"
	"sync"
)

// Map represents a threadsafe map of promises.
type Map struct {
	promises map[string]*Promise
	mutex    sync.Mutex
}

// NewMap instantiates and returns a new map.
func NewMap() *Map {
	return &Map{
		promises: make(map[string]*Promise),
		mutex:    sync.Mutex{},
	}
}

// Get returns the promise under the provided key.
func (m *Map) Get(key string) (*Promise, bool) {
	m.mutex.Lock()
	defer runtime.Gosched()
	defer m.mutex.Unlock()
	p, ok := m.promises[key]
	return p, ok
}

// Set will store a promise into the map under the provided key.
func (m *Map) Set(key string, p *Promise) {
	m.mutex.Lock()
	m.promises[key] = p
	m.mutex.Unlock()
	runtime.Gosched()
}

// Remove will remove a promise from the map.
func (m *Map) Remove(key string) {
	m.mutex.Lock()
	delete(m.promises, key)
	m.mutex.Unlock()
	runtime.Gosched()
}
