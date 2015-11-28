package es

import (
	"bufio"
	"compress/gzip"
	"io"
	"strings"

	"gopkg.in/olivere/elastic.v3"

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

func newReadyChan() chan error {
	// create and return a ready channel which is used to determine when ES is
	// ready to receive another bulk request
	rdy := make(chan error)
	go func() {
		rdy <- nil
	}()
	return rdy
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
	// get the ready channel
	rdy := newReadyChan()
	defer close(rdy)

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

	// scan file line by line
	scanner := bufio.NewScanner(reader)

	for {
		// create a new bulk request object
		bulk, err := GetBulkRequest(config.EsHost, config.EsPort, config.EsIndex, document.GetType())
		if err != nil {
			return err
		}
		for scanner.Scan() {
			// read line of file
			line := strings.Split(scanner.Text(), "\t")
			// set data for document
			document.SetData(line[0:])
			// get source from document
			source, err := document.GetSource()
			if err != nil {
				return err
			}
			// create index request
			// index := elastic.NewBulkIndexRequest().
			// 	Id(document.GetID()).
			// 	Doc(source)
			// size := index.String()
			// add index action to bulk req
			bulk.Add(
				elastic.NewBulkIndexRequest().
					Id(document.GetID()).
					Doc(source))

			if bulk.NumberOfActions()%config.BatchSize == 0 {
				// ready to send
				break
			}
		}

		// if no actions, we are finished
		if bulk.NumberOfActions() == 0 {
			break
		}

		// wait until last async request gets a response
		err = <-rdy
		if err != nil {
			return err
		}

		// send bulk request asynchronously
		go func(b *elastic.BulkService) {
			// send bulk ingest request
			rdy <- SendBulkRequest(b)
		}(bulk)
	}
	return <-rdy
}
