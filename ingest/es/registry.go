package es

import (
	"fmt"
)

var (
	// registry contains all document implementations for twitter data.
	registry = make(map[string]DocumentConstructor)
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

// DocumentConstructor represents a function to instantiate and return a
// document.
type DocumentConstructor func() Document

// Register registers a document constructor under the provided type id string.
func Register(typeID string, ctor DocumentConstructor) {
	registry[typeID] = ctor
}

// GetDocument when given a document id will return the document struct type.
func GetDocument(typeID string) (Document, error) {
	ctor, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Document type '%s' has not been registered", typeID)
	}
	return ctor(), nil
}
