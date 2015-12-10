package conf

import (
	"errors"
	"flag"
	"time"
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

// ParseCommandLine parses the commandline arguments and returns a Conf object.
func ParseCommandLine() (*Conf, error) {

	esHost := flag.String("es-host", "", "Elasticsearch host")
	esPort := flag.String("es-port", "9200", "Elasticsearch port")
	esIndex := flag.String("es-index", "", "Elasticsearch index")
	esDocType := flag.String("es-doc-type", "", "Elasticsearch type")
	esBatchSize := flag.Int("es-batch-size", 40000, "The bulk batch size in documents")
	esClearExisting := flag.Bool("es-clear-existing", true, "Clear index before ingest")
	hdfsHost := flag.String("hdfs-host", "", "HDFS host")
	hdfsPort := flag.String("hdfs-port", "", "HDFS port")
	hdfsPath := flag.String("hdfs-path", "", "HDFS ingest source data path")
	hdfsCompression := flag.String("hdfs-compression", "", "HDFS file compression used")
	startDate := flag.Int64("start-date", -1, "The unix timestamp (seconds) of the start date to ingest from")
	endDate := flag.Int64("end-date", -1, "The unix timestamp (seconds) of the end date to ingest to")
	duration := flag.Int64("duration", -1, "The duration in seconds to ingest either from start date, or end date, depending on which is provided")
	poolSize := flag.Int("pool-size", 8, "The worker pool size")

	// parse the flags
	flag.Parse()

	// check required flags
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

	// set dates
	start, end := parseTimeframe(*startDate, *endDate, *duration)

	// Set and save config
	config := &Conf{
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
		StartDate:       start,
		EndDate:         end,
		PoolSize:        *poolSize,
	}
	SaveConf(config)
	return config, nil
}
