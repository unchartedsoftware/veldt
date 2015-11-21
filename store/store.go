package store

type setInStore func( key string, value []byte ) error
type getFromStore func( key string ) ( []byte, error )
type existsInStore func( key string ) ( bool, error )

// Get returns a value from the store for a given string key.
var Get getFromStore = RedisGet

// Set will store a byte slice under a given key in the store.
var Set setInStore = RedisSet

// Exists returns whether or not a key exists in the store.
var Exists existsInStore = RedisExists
