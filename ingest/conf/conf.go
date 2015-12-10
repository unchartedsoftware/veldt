package conf

import (
	"time"
)

// Conf represents all the ingest runtime flags passed to the binary.
type Conf struct {
	// elasticsearch config
	EsHost          string
	EsPort          string
	EsIndex         string
	EsDocType       string
	EsClearExisting bool
	EsBatchSize     int
	// hdfs config
	HdfsHost        string
	HdfsPort        string
	HdfsPath        string
	HdfsCompression string
	// time filtering
	StartDate *time.Time
	EndDate   *time.Time
	// other
	PoolSize    int
}

var config *Conf

// SaveConf saves the parsed conf.
func SaveConf(c *Conf) {
	config = c
}

// GetConf returns a copy of the parsed conf.
func GetConf() Conf {
	return *config
}
