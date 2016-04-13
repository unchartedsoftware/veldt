package freecache

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/coocood/freecache"
	log "github.com/unchartedsoftware/plog"

	"github.com/unchartedsoftware/prism/store"
)

var (
	mutex = sync.Mutex{}
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
		cache.Set([]byte("derp:derp:derp"), []byte("test_data"), 0)
		got, err := cache.Get([]byte("derp:derp:derp"))
		if err != nil {
		    fmt.Println(err)
		} else {
		    fmt.Println(string(got))
		}
		caches[byteSize] = cache
	}
	mutex.Unlock()
	runtime.Gosched()
	return cache
}

// NewConnection instantiates and returns a new freecache store connection.
func NewConnection(byteSize int, expirySeconds int) store.ConnectionConstructor {
	return func() (store.Connection, error) {
		return &Connection{
			cache: getCache(byteSize),
			expiry: expirySeconds,
		}, nil
	}
}

// Get when given a string key will return a byte slice of data from redis.
func (r *Connection) Get(key string) ([]byte, error) {
	fmt.Printf("Get key: %s\n", key)
	fmt.Printf("Get key: %v\n", []byte(key))
 	return r.cache.Get([]byte(key))
}

// Set will store a byte slice under a given key in redis.
func (r *Connection) Set(key string, value []byte) error {
	fmt.Printf("Set key: %s\n", key)
	fmt.Printf("Set key: %v\n", []byte(key))
	return r.cache.Set([]byte(key), value, r.expiry)
}

// Exists returns whether or not a key exists in redis.
func (r *Connection) Exists(key string) (bool, error) {
	fmt.Printf("Exists key: %s\n", key)
	fmt.Printf("Exists key: %v\n", []byte(key))
	_, err := r.cache.Get([]byte(key))
	if err != nil {
		fmt.Printf("Exists err: %v\n", err)
		return false, nil
	}
	return true, nil
}

// Close closes the redis connection.
func (r *Connection) Close() {
	r.cache.Clear()
}
