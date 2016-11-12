package elastic

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/unchartedsoftware/plog"
	"gopkg.in/olivere/elastic.v3"
)

const (
	timeout = time.Second * 60
)

var (
	mutex   = sync.Mutex{}
	clients = make(map[string]*elastic.Client)
)

// NewClient instantiates and returns a new elastic.Client from a pool.
func NewClient(host string, port string) (*elastic.Client, error) {
	endpoint := host + ":" + port
	mutex.Lock()
	client, ok := clients[endpoint]
	if !ok {
		log.Infof("Connecting to elasticsearch `%s`", endpoint)
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
