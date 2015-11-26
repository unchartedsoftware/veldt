package es

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"strings"

	"github.com/unchartedsoftware/prism/ingest/conf"
	"github.com/unchartedsoftware/prism/ingest/hdfs"
	"github.com/unchartedsoftware/prism/ingest/info"
	"github.com/unchartedsoftware/prism/util/log"
)

// Index represents the 'index' property of an index action.
type Index struct {
	ID string `json:"_id"`
}

// IndexAction represents the 'index' action type for a bulk ingest.
type IndexAction struct {
	Index *Index `json:"index"`
}

// Document represents all necessary info to create an index and ingest a document.
type Document interface {
	Setup() error
	Teardown() error
	SetData([]string)
	GetSource() ([]byte, error)
	GetID() string
	GetMappings() string
	GetType() string
}

func getIndexAction(document Document) (*string, error) {
	// get source
	sourceBytes, sourceErr := document.GetSource()
	if sourceErr != nil {
		return nil, sourceErr
	}
	// build index action
	indexBytes, indexErr := json.Marshal(IndexAction{
		Index: &Index{
			ID: document.GetID(),
		},
	})
	if indexErr != nil {
		return nil, indexErr
	}
	jsonString := string(indexBytes) + "\n" + string(sourceBytes)
	return &jsonString, nil
}

func getDecompressReader(reader io.Reader, compression string) (io.Reader, error) {
	// use compression based reader if specified
	switch compression {
	case "gzip":
		return gzip.NewReader(reader)
	default:
		return reader, nil
	}
}

// IngestWorker is a worker to ingest twitter data into elasticsearch.
func IngestWorker(file info.IngestFile) error {
	// get the config struct
	config := conf.GetConf()
	// get document struct by type string
	document, err := GetDocumentByType(config.EsDocType)
	if err != nil {
		return err
	}
	// get data type
	dataType := document.GetType()

	// get hdfs file reader
	hdfsReader, err := hdfs.Open(config.HdfsHost, config.HdfsPort, file.Path+"/"+file.Name)
	if err != nil {
		return err
	}
	// defer close reader
	defer hdfsReader.Close()

	// get file reader
	reader, err := getDecompressReader(hdfsReader, config.HdfsCompression)
	if err != nil {
		return err
	}

	// create bulk actions slice
	actions := make([]string, config.BatchSize)
	actionIndex := 0

	// scan file line by line
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		// set data for document
		document.SetData(line[0:])
		action, err := getIndexAction(document)
		if err != nil {
			log.Error(err)
			continue
		}
		actions[actionIndex] = *action
		actionIndex++
		if actionIndex%config.BatchSize == 0 {
			// send bulk ingest request
			actionIndex = 0
			err := Bulk(config.EsHost, config.EsPort, config.EsIndex, dataType, actions[0:])
			if err != nil {
				return err
			}
		}
	}
	// send remaining documents
	if actionIndex > 0 {
		err = Bulk(config.EsHost, config.EsPort, config.EsIndex, dataType, actions[0:actionIndex])
		if err != nil {
			return err
		}
	}
	return nil
}
