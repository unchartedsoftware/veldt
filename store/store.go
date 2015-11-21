package store

type getFromStore func( key string ) ( []byte, error )
type setInStore func( key string, value []byte ) error

// Get returns a value from the store for a given string key.
var Get getFromStore = RedisGet

// Set will store a byte slice under a given key in the store.
var Set setInStore = RedisSet
