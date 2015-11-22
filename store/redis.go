package store

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	redisHost   = "localhost"
	redisPort   = "6379"
	maxIdle     = 8
	idleTimeout = 30 * time.Second
)

var redisPool = getPool(redisHost + ":" + redisPort)

func getPool(server string) *redis.Pool {
	fmt.Printf("Connecting to redis server: %s:%s\n", redisHost, redisPort)
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

func getConnection() redis.Conn {
	return redisPool.Get()
}

// RedisGet when given a string key will return a byte slice of data from redis.
func RedisGet(key string) ([]byte, error) {
	conn := getConnection()
	defer conn.Close()
	return redis.Bytes(conn.Do("GET", key))
}

// RedisSet will store a byte slice under a given key in redis.
func RedisSet(key string, value []byte) error {
	conn := getConnection()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

// RedisExists returns whether or not a key exists in redis.
func RedisExists(key string) (bool, error) {
	conn := getConnection()
	defer conn.Close()
	return redis.Bool(conn.Do("Exists", key))
}
