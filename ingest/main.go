package main

import (
    "flag"
    "fmt"
    "os"
	"runtime"
    "runtime/debug"

    "github.com/unchartedsoftware/prism/ingest/conf"
    "github.com/unchartedsoftware/prism/ingest/es"
    "github.com/unchartedsoftware/prism/ingest/info"
    "github.com/unchartedsoftware/prism/ingest/pool"
    "github.com/unchartedsoftware/prism/ingest/progress"
    "github.com/unchartedsoftware/prism/ingest/twitter"
)

var (
	esHost = flag.CommandLine.String( "es-host", "", "Elasticsearch host" )
    esPort = flag.CommandLine.String( "es-port", "9200", "Elasticsearch port" )
	esIndex = flag.CommandLine.String( "es-index", "", "Elasticsearch index" )
	esClearExisting = flag.CommandLine.Bool( "es-clear-existing", true, "Clear index before ingest" )
    hdfsHost = flag.CommandLine.String( "hdfs-host", "", "HDFS host" )
    hdfsPort = flag.CommandLine.String( "hdfs-port", "", "HDFS port" )
    hdfsPath = flag.CommandLine.String( "hdfs-path", "", "HDFS ingest source data path" )
    batchSize = flag.CommandLine.Int( "batch-size", 16000, "The bulk batch size in documents" )
    poolSize = flag.CommandLine.Int( "pool-size", 4, "The worker pool size" )
)

func parseArgs() conf.Conf {
    if *esHost == "" {
        fmt.Println("ElasticSearch host is not specified, please provide CL arg '-es-host'.")
        os.Exit(1)
    }
    if *esIndex == "" {
        fmt.Println("ElasticSearch index is not specified, please provide CL arg '-es-index'.")
        os.Exit(1)
    }
    if *hdfsHost == "" {
        fmt.Println("HDFS host is not specified, please provide CL arg '-hdfs-host'.")
        os.Exit(1)
    }
    if *hdfsPort == "" {
        fmt.Println("HDFS port is not specified, please provide CL arg '-hdfs-port'.")
        os.Exit(1)
    }
    if *hdfsPath == "" {
        fmt.Println("HDFS path is not specified, please provide CL arg '-hdfs-path'.")
        os.Exit(1)
    }
    config := conf.Conf{
        EsHost: *esHost,
        EsPort: *esPort,
        EsIndex: *esIndex,
        EsClearExisting: *esClearExisting,
        HdfsHost: *hdfsHost,
        HdfsPort: *hdfsPort,
        HdfsPath: *hdfsPath,
        BatchSize: *batchSize,
        PoolSize: *poolSize,
    }
	conf.SaveConf( &config )
    return config
}

func main() {

    runtime.GOMAXPROCS( runtime.NumCPU() )

    // parse flags
    flag.Parse()

    // load args into configig struct
    config := parseArgs()

    // check if index exists
    indexExists, err := es.IndexExists( config.EsHost, config.EsPort, config.EsIndex )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }

    // if index exists
    if indexExists && config.EsClearExisting {
        err = es.ClearIndex( config.EsHost, config.EsPort, config.EsIndex )
        if err != nil {
            fmt.Println("Failed to delete index, aborting.")
            debug.PrintStack()
            os.Exit(1)
        }
    }

    // if index does not exist at this point
    if !indexExists || config.EsClearExisting {
        err = es.CreateIndex(
            config.EsHost,
            config.EsPort,
            config.EsIndex,
            `{
                "mappings": ` + Twitter.Mappings + `
            }`)
        if err != nil {
            fmt.Println( err )
            debug.PrintStack()
            os.Exit(1)
        }
    }

    // get ingest info
    ingestInfo := info.GetIngestInfo( config.HdfsHost, config.HdfsPort, config.HdfsPath )

    // create pool of size N
    pool := pool.New( config.PoolSize )

    // display some info of the pending ingest
    fmt.Printf( "Indexing %d files containing %d bytes of data\n",
        len( ingestInfo.Files ),
        ingestInfo.NumTotalBytes )

    // launch the ingest job
    pool.Execute( twitter.TwitterWorker, ingestInfo )

    // finished succesfully
    progress.PrintTotalDuration()
}
