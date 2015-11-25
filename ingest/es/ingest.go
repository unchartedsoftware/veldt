package es

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/unchartedsoftware/prism/ingest/conf"
	"github.com/unchartedsoftware/prism/ingest/hdfs"
	"github.com/unchartedsoftware/prism/ingest/info"
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

func getFileReader(host string, port string, filepath string, compression string) (io.Reader, error) {
	// get hdfs file reader
	reader, err := hdfs.Open(host, port, filepath)
	if err != nil {
		return nil, err
	}
	// use compression based reader if specified
	switch compression {
	case "gzip":
		return gzip.NewReader(reader)
	default:
		return reader, nil
	}
}

// IngestWorker is a worker to ingest twitter data into elasticsearch.
func IngestWorker(file info.IngestFile) {
	config := conf.GetConf()
	documentTypeID := config.EsDocType
	indexType := GetDocumentByType(documentTypeID).GetType()
	documents := make([]string, config.BatchSize)
	documentIndex := 0

	// get file reader
	reader, err := getFileReader(config.HdfsHost, config.HdfsPort, file.Path+"/"+file.Name, config.HdfsCompression)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		// document := twitter.NYCTweetDocument{
		// 	Cols: line,
		// }
		document := GetDocumentByType(documentTypeID)
		document.SetData(line)
		action, err := getIndexAction(document)
		if err != nil {
			fmt.Println(err)
			debug.PrintStack()
			continue
		}
		documents[documentIndex] = *action
		documentIndex++
		if documentIndex%config.BatchSize == 0 {
			// send bulk ingest request
			documentIndex = 0
			err := Bulk(config.EsHost, config.EsPort, config.EsIndex, indexType, documents[0:])
			if err != nil {
				fmt.Println(err)
				debug.PrintStack()
				continue
			}
		}
	}
	//reader.Close()

	// send remaining documents
	err = Bulk(config.EsHost, config.EsPort, config.EsIndex, indexType, documents[0:documentIndex])
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
	}
}
