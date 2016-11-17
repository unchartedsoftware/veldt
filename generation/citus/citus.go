package citus

import (
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx"
	log "github.com/unchartedsoftware/plog"
)

const (
	timeout = time.Second * 60
)

var (
	mutex   = sync.Mutex{}
	clients = make(map[string]*pgx.ConnPool)
)

// NewClient return a citus client from the pool.
func NewClient(host string, port string) (*pgx.ConnPool, error) {
	endpoint := host + ":" + port
	portInt, err := strconv.Atoi(port)
	if err != nil {
		mutex.Unlock()
		runtime.Gosched()
		return nil, err
	}
	mutex.Lock()
	client, ok := clients[endpoint]
	if !ok {
		//TODO: Add configuration for connection parameters.
		log.Infof("Connecting to citus `%s`", endpoint)
		dbConfig := pgx.ConnConfig{
			Host: host,
			Port: uint16(portInt),
			User: "postgres",
		}

		poolConfig := pgx.ConnPoolConfig{
			ConnConfig:     dbConfig,
			MaxConnections: 16,
		}
		//TODO: Need to close the pool eventually. Not sure how to hook that in.
		c, err := pgx.NewConnPool(poolConfig)
		if err != nil {
			mutex.Unlock()
			runtime.Gosched()
			return nil, err
		}
		clients[endpoint] = c
		client = c
	}
	mutex.Unlock()
	runtime.Gosched()
	return client, nil
}
