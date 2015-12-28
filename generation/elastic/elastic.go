package elastic

import (
	"runtime"
	"sync"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/log"
)

var (
	mutex   = sync.Mutex{}
	clients = make(map[string]*elastic.Client)
)

func getClient(endpoint string) (*elastic.Client, error) {
	mutex.Lock()
	client, ok := clients[endpoint]
	if !ok {
		log.Debugf("Connecting to elasticsearch '%s'", endpoint)
		c, err := elastic.NewClient(
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
