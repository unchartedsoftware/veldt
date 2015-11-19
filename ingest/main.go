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
    "github.com/unchartedsoftware/prism/ingest/terms"
    "github.com/unchartedsoftware/prism/ingest/twitter"

    "github.com/unchartedsoftware/prism/util"
)

var (
	esHost = flag.CommandLine.String( "es-host", "", "Elasticsearch host" )
    esPort = flag.CommandLine.String( "es-port", "9200", "Elasticsearch port" )
	esIndex = flag.CommandLine.String( "es-index", "", "Elasticsearch index" )
	esType = flag.CommandLine.String( "es-type", "datum", "Elasticsearch type" )
	esClearExisting = flag.CommandLine.Bool( "es-clear-existing", true, "Clear index before ingest" )
    hdfsHost = flag.CommandLine.String( "hdfs-host", "", "HDFS host" )
    hdfsPort = flag.CommandLine.String( "hdfs-port", "", "HDFS port" )
    hdfsPath = flag.CommandLine.String( "hdfs-path", "", "HDFS ingest source data path" )
    batchSize = flag.CommandLine.Int( "batch-size", 16000, "The bulk batch size in documents" )
    poolSize = flag.CommandLine.Int( "pool-size", 8, "The worker pool size" )
    numTopTerms = flag.CommandLine.Int( "num-top-terms", 200, "The number of top terms to store" )
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
        EsType: *esType,
        EsClearExisting: *esClearExisting,
        HdfsHost: *hdfsHost,
        HdfsPort: *hdfsPort,
        HdfsPath: *hdfsPath,
        BatchSize: *batchSize,
        PoolSize: *poolSize,
        NumTopTerms: *numTopTerms,
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
        err = es.DeleteIndex( config.EsHost, config.EsPort, config.EsIndex )
        if err != nil {
            fmt.Println( err )
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
                "mappings": ` + twitter.GetMappings() + `
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
    fmt.Printf( "Processing %d files containing " + util.FormatBytes( float64( ingestInfo.NumTotalBytes ) ) + " of data\n",
        len( ingestInfo.Files ) )

    fmt.Println( "Determining top terms found in text" )

    // launch the top terms job
    pool.Execute( twitter.TwitterTopTermsWorker, ingestInfo )

    // finished succesfully
    progress.PrintTotalDuration()

    // save n current top term counts
    terms.SaveTopTerms( uint64( config.NumTopTerms ) );

    fmt.Println( "Ingesting data into elasticsearch" )

    // launch the ingest job
    pool.Execute( twitter.TwitterWorker, ingestInfo )

    // finished succesfully
    progress.PrintTotalDuration()
}
