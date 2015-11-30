package main

import (
	"errors"
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/unchartedsoftware/prism/ingest/conf"
	"github.com/unchartedsoftware/prism/ingest/es"
	"github.com/unchartedsoftware/prism/ingest/info"

	"github.com/unchartedsoftware/prism/util"
	"github.com/unchartedsoftware/prism/util/log"
)

var (
	esHost          = flag.String("es-host", "", "Elasticsearch host")
	esPort          = flag.String("es-port", "9200", "Elasticsearch port")
	esIndex         = flag.String("es-index", "", "Elasticsearch index")
	esDocType       = flag.String("es-doc-type", "", "Elasticsearch type")
	esBatchSize     = flag.Int("es-batch-size", 40000, "The bulk batch size in documents")
	esClearExisting = flag.Bool("es-clear-existing", true, "Clear index before ingest")
	hdfsHost        = flag.String("hdfs-host", "", "HDFS host")
	hdfsPort        = flag.String("hdfs-port", "", "HDFS port")
	hdfsPath        = flag.String("hdfs-path", "", "HDFS ingest source data path")
	hdfsCompression = flag.String("hdfs-compression", "", "HDFS file compression used")
	startDate       = flag.Int64("start-date", -1, "The unix timestamp (seconds) of the start date to ingest from")
	endDate         = flag.Int64("end-date", -1, "The unix timestamp (seconds) of the end date to ingest to")
	duration        = flag.Int64("duration", -1, "The duration in seconds to ingest either from start date, or end date, depending on which is provided")
	poolSize        = flag.Int("pool-size", 16, "The worker pool size")
	numTopTerms     = flag.Int("num-top-terms", 200, "The number of top terms to store")
)

func parseTimeframe(start int64, end int64, duration int64) (*time.Time, *time.Time) {
	var startDate time.Time
	var endDate time.Time
	// Viable options include and are parsed in order:
	//	1) startDate and endDate
	//	2) startDate and duration
	//	3) startDate only, endDate defaults to current time
	//  4) duration and endDate
	//	5) duration only, endDate defaults to current time
	if start != -1 {
		if end != -1 {
			// 1) startDate and endDate
			startDate = time.Unix(start, 0)
			endDate = time.Unix(end, 0)
		} else if duration != -1 {
			// 2) startDate and duration
			startDate = time.Unix(start, 0)
			endDate = startDate.Add(time.Second * time.Duration(duration))
		} else {
			// 3) startDate only, endDate defaults to current time
			startDate = time.Unix(start, 0)
			endDate = time.Now()
		}
		return &startDate, &endDate
	} else if duration != -1 {
		if end != -1 {
			// 4) endDate and duration
			endDate = time.Unix(end, 0)
			startDate = endDate.Add(time.Second * -time.Duration(duration))
		} else {
			// 5) duration only, endDate defaults to current time
			endDate = time.Now()
			startDate = endDate.Add(time.Second * -time.Duration(duration))
		}
		return &startDate, &endDate
	}
	return nil, nil
}

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
	startDate, endDate := parseTimeframe(*startDate, *endDate, *duration)
	config := &conf.Conf{
		EsHost:          *esHost,
		EsPort:          *esPort,
		EsIndex:         *esIndex,
		EsDocType:       *esDocType,
		EsBatchSize:     *esBatchSize,
		EsClearExisting: *esClearExisting,
		HdfsHost:        *hdfsHost,
		HdfsPort:        *hdfsPort,
		HdfsPath:        *hdfsPath,
		HdfsCompression: *hdfsCompression,
		StartDate:       startDate,
		EndDate:         endDate,
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

	// load args into config struct
	config, err := parseArgs()
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
