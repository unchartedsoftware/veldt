package hdfs

import (
	"os"

	"github.com/colinmarc/hdfs"

	"github.com/unchartedsoftware/prism/util/log"
)

var hdfsClient *hdfs.Client

func getHdfsClient(host string, port string) (*hdfs.Client, error) {
	endpoint := host + ":" + port
	if hdfsClient == nil {
		log.Debug("Connecting to HDFS: " + endpoint)
		client, err := hdfs.New(endpoint)
		hdfsClient = client
		return hdfsClient, err
	}
	return hdfsClient, nil
}

// Open will connect to HDFS, open a file, and return a file reader.
func Open(host string, port string, filePath string) (*hdfs.FileReader, error) {
	// get hdfs client
	client, err := getHdfsClient(host, port)
	if err != nil {
		return nil, err
	}
	// get hdfs file reader
	return client.Open(filePath)
}

// ReadDir will connect to HDFS and return an array containing file information.
func ReadDir(host string, port string, path string) ([]os.FileInfo, error) {
	// get hdfs client
	client, err := getHdfsClient(host, port)
	if err != nil {
		return nil, err
	}
	return client.ReadDir(path)
}
