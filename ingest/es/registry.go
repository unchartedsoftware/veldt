package es

import (
	"fmt"
	"os"
	"reflect"

	"github.com/unchartedsoftware/prism/ingest/twitter"
)

// registry contains all document implementations for twitter data.
var registry = make(map[string]reflect.Type)

// register all document implementations here.
func init() {
	registry["nyc_twitter"] = reflect.TypeOf(twitter.NYCTweetDocument{})
	registry["isil_twitter"] = reflect.TypeOf(twitter.ISILTweetDocument{})
}

// GetDocumentByType when given a document id will return the document struct type.
func GetDocumentByType(typeID string) Document {
	docType, ok := registry[typeID]
	if ok {
		return reflect.New(docType).Interface().(Document)
	}
	fmt.Println("Document type '" + typeID + "' has not been registered.")
	defer os.Exit(1)
	return nil
}
