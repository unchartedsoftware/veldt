package redis

import (
	"runtime"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/unchartedsoftware/prism/log"
)

const (
	maxIdle     = 8
	idleTimeout = 10 * time.Second
)

var (
	mutex = sync.Mutex{}
	pools = make(map[string]*redis.Pool)
)

func getConnection(endpoint string) redis.Conn {
	mutex.Lock()
	pool, ok := pools[endpoint]
	if !ok {
		log.Debugf("Connecting to redis 'tcp://%s'", endpoint)
		pool = newPool(endpoint)
		pools[endpoint] = pool
	}
	mutex.Unlock()
	runtime.Gosched()
	return pool.Get()
}

func newPool(endpoint string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", endpoint)
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
