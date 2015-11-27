package main

import (
	"errors"
	"flag"
	"os"
	"runtime"

	"github.com/unchartedsoftware/prism/ingest/conf"
	"github.com/unchartedsoftware/prism/ingest/es"
	"github.com/unchartedsoftware/prism/ingest/info"
	"github.com/unchartedsoftware/prism/ingest/pool"
	"github.com/unchartedsoftware/prism/ingest/progress"

	"github.com/unchartedsoftware/prism/util"
	"github.com/unchartedsoftware/prism/util/log"
)

var (
	esHost          = flag.CommandLine.String("es-host", "", "Elasticsearch host")
	esPort          = flag.CommandLine.String("es-port", "9200", "Elasticsearch port")
	esIndex         = flag.CommandLine.String("es-index", "", "Elasticsearch index")
	esDocType       = flag.CommandLine.String("es-doc-type", "", "Elasticsearch type")
	esClearExisting = flag.CommandLine.Bool("es-clear-existing", true, "Clear index before ingest")
	hdfsHost        = flag.CommandLine.String("hdfs-host", "", "HDFS host")
	hdfsPort        = flag.CommandLine.String("hdfs-port", "", "HDFS port")
	hdfsPath        = flag.CommandLine.String("hdfs-path", "", "HDFS ingest source data path")
	hdfsCompression = flag.CommandLine.String("hdfs-compression", "", "HDFS file compression used")
	batchSize       = flag.CommandLine.Int("batch-size", 24000, "The bulk batch size in documents")
	poolSize        = flag.CommandLine.Int("pool-size", 8, "The worker pool size")
	numTopTerms     = flag.CommandLine.Int("num-top-terms", 200, "The number of top terms to store")
)

func parseArgs() (*conf.Conf, error) {
	if *esHost == "" {
		return nil, errors.New("ElasticSearch host is not specified, please provide CL arg '-es-host'")
	}
	if *esIndex == "" {
		return nil, errors.New("ElasticSearch index is not specified, please provide CL arg '-es-index'")
	}
	if *esDocType == "" {
		return nil, errors.New("ElasticSearch document type is not specified, please provide CL arg '-es-doc-type'")
	}
	if *hdfsHost == "" {
		return nil, errors.New("HDFS host is not specified, please provide CL arg '-hdfs-host'")
	}
	if *hdfsPort == "" {
		return nil, errors.New("HDFS port is not specified, please provide CL arg '-hdfs-port'")
	}
	if *hdfsPath == "" {
		return nil, errors.New("HDFS path is not specified, please provide CL arg '-hdfs-path'")
	}
	config := &conf.Conf{
		EsHost:          *esHost,
		EsPort:          *esPort,
		EsIndex:         *esIndex,
		EsDocType:       *esDocType,
		EsClearExisting: *esClearExisting,
		HdfsHost:        *hdfsHost,
		HdfsPort:        *hdfsPort,
		HdfsPath:        *hdfsPath,
		HdfsCompression: *hdfsCompression,
		BatchSize:       *batchSize,
		PoolSize:        *poolSize,
		NumTopTerms:     *numTopTerms,
	}
	conf.SaveConf(config)
	return config, nil
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse flags
	flag.Parse()

	// load args into configig struct
	config, err := parseArgs()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// check document type id and implementation
	document, err := es.GetDocumentByType(config.EsDocType)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// get ingest info
	ingestInfo, errs := info.GetIngestInfo(config.HdfsHost, config.HdfsPort, config.HdfsPath)
	if errs != nil {
		log.Error(errs)
		os.Exit(1)
	}

	// create pool of size N
	pool := pool.New(config.PoolSize)

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

	// launch the ingest job
	err = pool.Execute(es.IngestWorker, ingestInfo)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// teardown after ingest
	document.Teardown()

	// finished succesfully
	progress.PrintTotalDuration()
}
