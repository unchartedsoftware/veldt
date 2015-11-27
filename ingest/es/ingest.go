package es

import (
	"bufio"
	"compress/gzip"
	"io"
	"strings"

	"github.com/unchartedsoftware/prism/ingest/conf"
	"github.com/unchartedsoftware/prism/ingest/hdfs"
	"github.com/unchartedsoftware/prism/ingest/info"
)

// Document represents all necessary info to create an index and ingest a document.
type Document interface {
	Setup() error
	Teardown() error
	SetData([]string)
	GetSource() (interface{}, error)
	GetID() string
	GetMappings() string
	GetType() string
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
	documents := make([]*Document, config.BatchSize)
	documentIndex := 0

	// scan file line by line
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		// get document struct by type string
		document, err := GetDocumentByType(config.EsDocType)
		if err != nil {
			return err
		}
		// set data for document
		document.SetData(line[0:])
		documents[documentIndex] = &document
		documentIndex++
		if documentIndex%config.BatchSize == 0 {
			// send bulk ingest request
			documentIndex = 0
			_, err := Bulk(config.EsHost, config.EsPort, config.EsIndex, documents[0:])
			if err != nil {
				return err
			}
		}
	}
	// send remaining documents
	if documentIndex > 0 {
		_, err = Bulk(config.EsHost, config.EsPort, config.EsIndex, documents[0:documentIndex])
		if err != nil {
			return err
		}
	}
	return nil
}
