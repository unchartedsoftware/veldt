package twitter

import (
	"bufio"
	"os"
	"strings"

	"github.com/unchartedsoftware/prism/ingest/conf"
	"github.com/unchartedsoftware/prism/ingest/hdfs"
	"github.com/unchartedsoftware/prism/ingest/terms"
	"github.com/unchartedsoftware/prism/util/log"
)

// TopTermsWorker is a worker to calculate the top terms found in tweet text.
func TopTermsWorker(file os.FileInfo) {
	config := conf.GetConf()

	// get hdfs file reader
	reader, err := hdfs.Open(config.HdfsHost, config.HdfsPort, config.HdfsPath+"/"+file.Name())
	if err != nil {
		log.Error(err)
		return
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		terms.AddTerms(line[4])
	}

	reader.Close()
}
