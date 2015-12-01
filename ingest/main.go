package main

import (
	"os"
	"runtime"

	"github.com/unchartedsoftware/prism/ingest/conf"
	"github.com/unchartedsoftware/prism/ingest/es"
	"github.com/unchartedsoftware/prism/ingest/info"

	"github.com/unchartedsoftware/prism/util"
	"github.com/unchartedsoftware/prism/util/log"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse flags into config struct
	config, err := conf.ParseCommandLine()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// check document type and implementation
	document, err := es.GetDocumentByType(config.EsDocType)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// get ingest info
	ingestInfo, errs := info.GetIngestInfo(config.HdfsHost, config.HdfsPort, config.HdfsPath, document)
	if errs != nil {
		log.Error(errs)
		os.Exit(1)
	}

	// display some info of the pending ingest
	log.Debugf("Processing %d files containing "+util.FormatBytes(float64(ingestInfo.NumTotalBytes))+" of data",
		len(ingestInfo.Files))

	// prepare elasticsearch index
	err = es.PrepareIndex(config.EsHost, config.EsPort, config.EsIndex, document.GetMappings(), config.EsClearExisting)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// setup for ingest
	document.Setup()

	// create pool of size N
	pool := info.NewPool(config.PoolSize)
	// launch the ingest job
	err = pool.Execute(info.IngestWorker, ingestInfo)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// teardown after ingest
	document.Teardown()
}
