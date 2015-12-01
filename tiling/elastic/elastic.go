package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/log"
)

const (
	esHost  = "http://10.64.16.120:9200"
	esIndex = "nyc_twitter"
)

var (
	client = getClient()
)

func getClient() /*map[string]*/ *elastic.Client {
	log.Debugf("Connecting to elasticsearch '%s/%s'", esHost)
	client, err := elastic.NewClient(
		elastic.SetURL(esHost),
		elastic.SetSniff(false),
		elastic.SetGzip(true),
	)
	if err != nil {
		log.Error(err)
		return nil
	}
	return client
}
