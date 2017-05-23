package elastic

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
)

const (
	timeout = time.Second * 60
)

var (
	mutex   = sync.Mutex{}
	clients = make(map[string]*elastic.Client)
)

// Elastic represents an elasticsearch type.
type Elastic struct {
	Host string
	Port string
}

// CreateSearchService creates the elasticsearch search service from the provided uri.
func (e *Elastic) CreateSearchService(uri string) (*elastic.SearchService, error) {
	// get client
	client, err := e.createClient()
	if err != nil {
		return nil, err
	}
	split := strings.Split(uri, "/")
	if len(split) < 2 {
		index := split[0]
		return client.Search().
			Index(index).
			Size(0), nil
	}
	index := split[0]
	typ := split[1]
	return client.Search().
		Index(index).
		Type(typ).
		Size(0), nil
}

// CreateQuery creates the elasticsearch query from the query struct.
func (e *Elastic) CreateQuery(query veldt.Query) (*elastic.BoolQuery, error) {
	// create root query
	root := elastic.NewBoolQuery()
	// add filter query
	if query != nil {
		// type assert
		esquery, ok := query.(Query)
		if !ok {
			return nil, fmt.Errorf("query is not elastic.Query")
		}
		// get underlying query
		q, err := esquery.Get()
		if err != nil {
			return nil, err
		}
		root.Must(q)
	}
	return root, nil
}

// CreateMappingService creates the elasticsearch mapping service from the provided uri.
func (e *Elastic) CreateMappingService(uri string) (*elastic.IndicesGetMappingService, error) {
	// get client
	client, err := e.createClient()
	if err != nil {
		return nil, err
	}
	split := strings.Split(uri, "/")
	if len(split) < 2 {
		index := split[0]
		return client.GetMapping().
			Index(index), nil
	}
	index := split[0]
	typ := split[1]
	return client.GetMapping().
		Index(index).
		Type(typ), nil
}

func (e *Elastic) createClient() (*elastic.Client, error) {
	endpoint := e.Host + ":" + e.Port
	mutex.Lock()
	client, ok := clients[endpoint]
	if !ok {
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
