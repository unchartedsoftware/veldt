package es

import (
	"bufio"
	"compress/gzip"
	"io"
	"strings"
	"time"

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

func timestamp() uint64 {
	// return timestamp in ms
    return uint64( time.Now().UnixNano() ) / uint64(time.Millisecond)
}

// IngestWorker is a worker to ingest twitter data into elasticsearch.
func IngestWorker(file info.IngestFile,  eq *Equalizer) error {

	// get the config struct
	config := conf.GetConf()
	// get document struct by type string
	document, err := GetDocumentByType(config.EsDocType)
	if err != nil {
		return err
	}

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

		// get current timestamp, this will be used to calculate how long it took
		// to generate the bulk payload
		ts := timestamp()

		// begin reading file, line by line
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

		// wait until the equalizer determines ES is ready, also check the
		// status of the last req, if error, return error
		rErr := <- eq.Ready
		if rErr != nil {
			return rErr
		}

		// send bulk request asynchronously
		eq.Send <- Request {
			Bulk: bulk,
			Took: timestamp() - ts, // how long it took to generate payload
		}
	}
	return nil
}
