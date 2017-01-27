package redis

import (
	"github.com/garyburd/redigo/redis"

	"github.com/unchartedsoftware/veldt"
)

// Store represents a single connection to a redis server.
type Store struct {
	conn   redis.Conn
	expiry int
}

// NewStore instantiates and returns a new redis store connection.
func NewStore(host, port string, expirySeconds int) veldt.StoreCtor {
	return func() (veldt.Store, error) {
		return &Store{
			conn:   getConnection(host, port),
			expiry: expirySeconds,
		}, nil
	}
}

// Get when given a string key will return a byte slice of data from redis.
func (r *Store) Get(key string) ([]byte, error) {
	return redis.Bytes(r.conn.Do("GET", key))
}

// Set will store a byte slice under a given key in redis.
func (r *Store) Set(key string, value []byte) error {
	var err error
	if r.expiry > 0 {
		_, err = r.conn.Do("SET", key, value, "NX", "EX", r.expiry)
	} else {
		_, err = r.conn.Do("SET", key, value)
	}
	return err
}

// Exists returns whether or not a key exists in redis.
func (r *Store) Exists(key string) (bool, error) {
	return redis.Bool(r.conn.Do("Exists", key))
}

// Close closes the redis connection.
func (r *Store) Close() {
	r.conn.Close()
}
