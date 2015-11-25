package es

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// Bulk sends a bulk request to elasticsearch with the provided payload.
func Bulk(host string, port string, index string, datatype string, actions []string) error {
	jsonLines := fmt.Sprintf("%s\n", strings.Join(actions, "\n"))
	response, err := http.Post(host+":"+port+"/"+index+"/"+datatype+"/_bulk", "application/json", strings.NewReader(jsonLines))
	if err != nil {
		fmt.Println(err)
		return errors.New("Bulk request failed")
	}
	response.Body.Close()
	return nil
}

// IndexExists returns whether or not the provided index exists in elasticsearch.
func IndexExists(host string, port string, index string) (bool, error) {
	resp, _, errs := gorequest.New().
		Head(host + ":" + port + "/" + index).
		End()
	if errs != nil {
		fmt.Println(errs)
		return false, errors.New("Unable to determine if index exists")
	}
	return resp.StatusCode != 404, nil
}

// DeleteIndex deletes the provided index in elasticsearch.
func DeleteIndex(host string, port string, index string) error {
	fmt.Println("Clearing index '" + index + "'")
	_, _, errs := gorequest.New().
		Delete(host + ":" + port + "/" + index).
		End()
	if errs != nil {
		fmt.Println(errs)
		return errors.New("Failed to delete index")
	}
	return nil
}

// CreateIndex creates the provided index in elasticsearch.
func CreateIndex(host string, port string, index string, body string) error {
	fmt.Println("Creating index '" + index + "'")
	_, _, errs := gorequest.New().
		Put(host + ":" + port + "/" + index).
		Send(body).
		End()
	if errs != nil {
		fmt.Println(errs)
		return errors.New("Failed to create index")
	}
	return nil
}

// PrepareIndex will ensure the provided index exists, and will optionally clear it.
func PrepareIndex(host string, port string, index string, documentTypeID string, clearExisting bool) error {
	// check if index exists
	indexExists, err := IndexExists(host, port, index)
	if err != nil {
		return err
	}
	// if index exists
	if indexExists && clearExisting {
		err = DeleteIndex(host, port, index)
		if err != nil {
			return err
		}
	}
	// get document struct by type id
	mappings := GetDocumentByType(documentTypeID).GetMappings()
	// if index does not exist at this point
	if !indexExists || clearExisting {
		err = CreateIndex(
			host,
			port,
			index,
			`{
                "mappings": `+mappings+`
            }`)
		if err != nil {
			return err
		}
	}
	return nil
}
