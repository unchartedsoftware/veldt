package store

import (
	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool
var redisHost = "localhost"
var redisPort = "6379"
var maxConnections = 64

func getConnection() redis.Conn {
    if redisPool == nil {
        redisPool = redis.NewPool( func() ( redis.Conn, error ) {
    		conn, err := redis.Dial( "tcp", redisHost + ":" + redisPort )
    		if err != nil {
    			return nil, err
    		}
    		return conn, err
    	}, maxConnections )
    }
    return redisPool.Get()
}

// RedisGet when given a string key will return a byte slice of data from redis.
func RedisGet( key string ) ( []byte, error ) {
    conn := getConnection()
    defer conn.Close()
    return redis.Bytes( conn.Do( "GET", key ) )
}

// RedisSet will store a byte slice under a given key in redis.
func RedisSet( key string, value []byte ) error {
    conn := getConnection()
    defer conn.Close()
    _, err := conn.Do( "SET", key, value )
	return err
}

// RedisExists returns whether or not a key exists in redis.
func RedisExists( key string ) ( bool, error ) {
    conn := getConnection()
    defer conn.Close()
    return redis.Bool( conn.Do( "Exists", key ) )
}
