package es

import (
	"errors"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/log"
)

var esClient *elastic.Client

func getClient(host string, port string) (*elastic.Client, error) {
	endpoint := host + ":" + port
	if esClient == nil {
		client, err := elastic.NewClient(
			elastic.SetURL(endpoint),
			elastic.SetSniff(false),
			elastic.SetGzip(true),
		)
		if err != nil {
			return nil, err
		}
		esClient = client
	}
	return esClient, nil
}

// GetBulkRequest creates and returns a pointer to a new elastic.BulkService
// for building a bulk request.
func GetBulkRequest(host string, port string, index string, typ string) (*elastic.BulkService, error) {
	client, err := getClient(host, port)
	if err != nil {
		return nil, err
	}
	return client.Bulk().
		Index(index).
		Type(typ), nil
}

// NewBulkIndexRequest creates and returns a pointer to a BulkIndexRequest object.
func NewBulkIndexRequest() *elastic.BulkIndexRequest {
	return elastic.NewBulkIndexRequest()
}

// SendBulkRequest sends the provided bulk request and handles the response.
func SendBulkRequest(bulk *elastic.BulkService) (uint64, error) {
	res, err := bulk.Do()
	if err != nil {
		return 0, err
	}
	if res.Errors {
		// find first error and return it
		for _, item := range res.Items {
			if item["index"].Error != nil {
				return uint64(res.Took), fmt.Errorf("%s, %s", item["index"].Error.Type, item["index"].Error.Reason)
			}
		}
	}
	return uint64(res.Took), nil
}

// IndexExists returns whether or not the provided index exists in elasticsearch.
func IndexExists(host string, port string, index string) (bool, error) {
	client, err := getClient(host, port)
	if err != nil {
		return false, err
	}
	return client.IndexExists(index).Do()
}

// DeleteIndex deletes the provided index in elasticsearch.
func DeleteIndex(host string, port string, index string) error {
	client, err := getClient(host, port)
	if err != nil {
		return err
	}
	res, err := client.DeleteIndex(index).Do()
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return errors.New("Delete index request not acknowledged")
	}
	return nil
}

// CreateIndex creates the provided index in elasticsearch.
func CreateIndex(host string, port string, index string, body string) error {
	client, err := getClient(host, port)
	if err != nil {
		return err
	}
	res, err := client.CreateIndex(index).Body(body).Do()
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return errors.New("Create index request not acknowledged")
	}
	return nil
}

// PrepareIndex will ensure the provided index exists, and will optionally clear it.
func PrepareIndex(host string, port string, index string, mappings string, clearExisting bool) error {
	// check if index exists
	indexExists, err := IndexExists(host, port, index)
	if err != nil {
		return err
	}
	// if index exists
	if indexExists && clearExisting {
		log.Debug("Deleting index '" + index + "'")
		err = DeleteIndex(host, port, index)
		if err != nil {
			return err
		}
	}
	// if index does not exist at this point, create it
	if !indexExists || clearExisting {
		log.Debug("Creating index '" + index + "'")
		err = CreateIndex(host, port, index, `{"mappings":`+mappings+`}`)
		if err != nil {
			return err
		}
	}
	return nil
}
