package redis

import (
	"github.com/garyburd/redigo/redis"

	"github.com/unchartedsoftware/prism/store"
)

// Connection represents a single connection to a redis server.
type Connection struct {
	conn redis.Conn
}

// NewConnection instantiates and returns a new redis store connection.
func NewConnection(req *store.Request) (store.Connection, error) {
	return &Connection{
		conn: getConnection(req.Endpoint),
	}, nil
}

// Get when given a string key will return a byte slice of data from redis.
func (r *Connection) Get(key string) ([]byte, error) {
	return redis.Bytes(r.conn.Do("GET", key))
}

// Set will store a byte slice under a given key in redis.
func (r *Connection) Set(key string, value []byte) error {
	_, err := r.conn.Do("SET", key, value)
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
