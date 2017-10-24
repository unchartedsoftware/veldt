package citus

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/jackc/pgx"
)

const (
	timeout = time.Second * 60
)

var (
	mutex   = sync.Mutex{}
	clients = make(map[string]*pgx.ConnPool)
)

// Config defines the database details required to establish a connection.
type Config struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

// NewClient return a citus client from the pool.
func NewClient(cfg *Config) (*pgx.ConnPool, error) {
	endpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	mutex.Lock()
	client, ok := clients[endpoint]
	if !ok {
		//TODO: Add configuration for connection parameters.
		dbConfig := pgx.ConnConfig{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Database: cfg.Database,
			User:     cfg.User,
			Password: cfg.Password,
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
