package hdfs

import (
	"os"

	"github.com/colinmarc/hdfs"

	"github.com/unchartedsoftware/prism/log"
)

const (
	hdfsUser = "etl"
)

var (
	hdfsClient *hdfs.Client
)

func getHdfsClient(host string, port string) (*hdfs.Client, error) {
	endpoint := host + ":" + port
	if hdfsClient == nil {
		log.Debug("Connecting to HDFS: " + endpoint)
		client, err := hdfs.NewForUser(endpoint, hdfsUser)
		if err != nil {
			return nil, err
		}
		hdfsClient = client
	}
	return hdfsClient, nil
}

// Open will connect to HDFS, open a file, and return a file reader.
func Open(host string, port string, filePath string) (*hdfs.FileReader, error) {
	client, err := getHdfsClient(host, port)
	if err != nil {
		return nil, err
	}
	return client.Open(filePath)
}

// ReadDir will connect to HDFS and return an array containing file information.
func ReadDir(host string, port string, path string) ([]os.FileInfo, error) {
	client, err := getHdfsClient(host, port)
	if err != nil {
		return nil, err
	}
	return client.ReadDir(path)
}
