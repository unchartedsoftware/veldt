package redis

import (
	"github.com/garyburd/redigo/redis"

	"github.com/unchartedsoftware/prism/store"
)

// Connection represents a single connection to a redis server.
type Connection struct {
	conn   redis.Conn
	expiry string
}

// NewConnection instantiates and returns a new redis store connection.
func NewConnection(host, port, expirySeconds string) store.ConnectionConstructor {
	return func() (store.Connection, error) {
		return &Connection{
			conn:   getConnection(host, port),
			expiry: expirySeconds,
		}, nil
	}
}

// Get when given a string key will return a byte slice of data from redis.
func (r *Connection) Get(key string) ([]byte, error) {
	return redis.Bytes(r.conn.Do("GET", key))
}

// Set will store a byte slice under a given key in redis.
func (r *Connection) Set(key string, value []byte) error {
	var err error
	if r.expiry != "" {
		_, err = r.conn.Do("SET", key, value, "NX", "EX", r.expiry)
	} else {
		_, err = r.conn.Do("SET", key, value)
	}
	return err
}

// Exists returns whether or not a key exists in redis.
func (r *Connection) Exists(key string) (bool, error) {
	return redis.Bool(r.conn.Do("Exists", key))
}

// Close closes the redis connection.
func (r *Connection) Close() {
	r.conn.Close()
}
