package es

import (
	"errors"
	"reflect"

	"github.com/unchartedsoftware/prism/ingest/twitter"
)

// Document represents all necessary info to create an index and ingest a document.
type Document interface {
	Setup() error
	Teardown() error
	FilterDir(string) bool
	FilterFile(string) bool
	SetData([]string)
	GetSource() (interface{}, error)
	GetID() string
	GetMappings() string
	GetType() string
}

// registry contains all document implementations for twitter data.
var registry = make(map[string]reflect.Type)

// register all document implementations here.
func init() {
	registry["nyc_twitter"] = reflect.TypeOf(twitter.NYCTweet{})
	registry["isil_twitter"] = reflect.TypeOf(twitter.ISILTweet{})
	registry["isil_twitter_deprecated"] = reflect.TypeOf(twitter.ISILTweetDeprecated{})
}

// GetDocumentByType when given a document id will return the document struct type.
func GetDocumentByType(typeID string) (Document, error) {
	docType, ok := registry[typeID]
	if !ok {
		return nil, errors.New("Document type '" + typeID + "' has not been registered.")
	}
	doc, ok := reflect.New(docType).Interface().(Document)
	if !ok {
		return nil, errors.New("Document type '" + typeID + "' does not implement the Document interface.")
	}
	return doc, nil
}
