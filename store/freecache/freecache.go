package freecache

import (
	"runtime"
	"sync"

	"github.com/coocood/freecache"
	"github.com/unchartedsoftware/plog"

	"github.com/unchartedsoftware/veldt"
)

var (
	mutex  = sync.Mutex{}
	caches = make(map[int]*freecache.Cache)
)

// Connection represents a single connection to a freecache instance.
type Connection struct {
	cache  *freecache.Cache
	expiry int
}

func getCache(byteSize int) *freecache.Cache {
	mutex.Lock()
	cache, ok := caches[byteSize]
	if !ok {
		log.Infof("Creating freecache instance of size `%d` bytes", byteSize)
		cache = freecache.NewCache(byteSize)
		caches[byteSize] = cache
	}
	mutex.Unlock()
	runtime.Gosched()
	return cache
}

// NewConnection instantiates and returns a new freecache store connection.
func NewConnection(byteSize int, expirySeconds int) veldt.StoreCtor {
	return func() (veldt.Store, error) {
		return &Connection{
			cache:  getCache(byteSize),
			expiry: expirySeconds,
		}, nil
	}
}

// Get when given a string key will return a byte slice of data from freecache.
func (r *Connection) Get(key string) ([]byte, error) {
	return r.cache.Get([]byte(key))
}

// Set will store a byte slice under a given key in freecache.
func (r *Connection) Set(key string, value []byte) error {
	return r.cache.Set([]byte(key), value, r.expiry)
}

// Exists returns whether or not a key exists in freecache.
func (r *Connection) Exists(key string) (bool, error) {
	_, err := r.cache.Get([]byte(key))
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Close closes the freecache connection.
func (r *Connection) Close() {
	// no-op
}
