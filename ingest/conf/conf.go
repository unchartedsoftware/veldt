package conf

type Conf struct {
    EsHost string
    EsPort string
    EsIndex string
    EsClearExisting bool
    HdfsHost string
    HdfsPort string
    HdfsPath string
    BatchSize int
    PoolSize int
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
