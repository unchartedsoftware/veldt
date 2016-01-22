package elastic

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	log "github.com/unchartedsoftware/plog"
	"gopkg.in/olivere/elastic.v3"
)

const (
	timeout = time.Second * 30
)

var (
	mutex   = sync.Mutex{}
	clients = make(map[string]*elastic.Client)
)

// GetClient returns an elasticsearch client from the pool.
func GetClient(endpoint string) (*elastic.Client, error) {
	mutex.Lock()
	client, ok := clients[endpoint]
	if !ok {
		log.Debugf("Connecting to elasticsearch '%s'", endpoint)
		c, err := elastic.NewClient(
			elastic.SetHttpClient(&http.Client{
				Timeout: timeout,
			}),
			elastic.SetURL(endpoint),
			elastic.SetSniff(false),
			elastic.SetGzip(true),
		)
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
