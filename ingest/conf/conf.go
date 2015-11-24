package conf

// Conf represents all the ingest runtime flags passed to the binary.
type Conf struct {
	EsHost          string
	EsPort          string
	EsIndex         string
	EsDocType       string
	EsClearExisting bool
	HdfsHost        string
	HdfsPort        string
	HdfsPath        string
	BatchSize       int
	PoolSize        int
	NumTopTerms     int
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
