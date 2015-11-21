package store

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool
var redisHost = "localhost"
var redisPort = "6379"
var maxConnections = 10

func getConnection() *redis.Pool {
    if redisPoll == nil {
        redisPool := redis.NewPool( func() ( redis.Conn, error ) {
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
    return redis.String( c.Do( "GET", key ) )
}

// RedisSet will store a byte slice under a given key in redis.
func RedisSet( key string, value []byte ) error {
    conn := getConnection()
    defer conn.Close()
    status, err := conn.Do( "SET", key, value )
}
