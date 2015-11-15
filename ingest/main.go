package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "runtime/debug"
    "strings"

    "github.com/unchartedsoftware/prism/ingest/es"
    "github.com/unchartedsoftware/prism/ingest/hdfs"
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
)

type IngestInfo struct {
    Paths []string
    FileSizes []int64
    NumTotalBytes int64
}

func getFileInfo() *IngestInfo {
    client, err := hdfs.GetHdfsClient( conf.HdfsHost, conf.HdfsPort )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }
    fmt.Println( "Retreiving ingest directory information from: " + conf.HdfsPath )
    files, err := client.ReadDir( conf.HdfsPath )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }
    var paths []string
    var fileSizes []int64
    var numTotalBytes int64 = 0
    for i:= 0; i<len( files ); i++ {
        file := files[i]
        if !file.IsDir() && file.Name() != ".SUCCESS" && file.Size() > 0 {
            // add to total bytes
            numTotalBytes += file.Size()
            // store path and file length
            paths = append( paths, conf.HdfsPath + "/" + file.Name() )
            fileSizes = append( fileSizes, file.Size() )
        }
    }
    return &IngestInfo{
        Paths: paths,
        FileSizes: fileSizes,
        NumTotalBytes: numTotalBytes,
    }
}

func ingestFile( filePath string ) {
    const batchSize = 9999
    documents := make( []string, batchSize )
    documentIndex := 0

    // get hdfs client
    client, clientErr := hdfs.GetHdfsClient( conf.HdfsHost, conf.HdfsPort )
    if clientErr != nil {
        fmt.Println( clientErr )
        debug.PrintStack()
        os.Exit(1)
    }

    reader, fileErr := client.Open( filePath )
    if fileErr != nil {
        fmt.Println( fileErr )
        debug.PrintStack()
        os.Exit(1)
    }

    scanner := bufio.NewScanner( reader )
    for scanner.Scan() {
        line := strings.Split( scanner.Text(), "\t" )
        action, err := twitter.CreateIndexAction( line )
        if err != nil {
            fmt.Println( err )
            debug.PrintStack()
            continue
        }
        documents[ documentIndex ] = *action
        documentIndex++
        if documentIndex % batchSize == 0 {
            // send bulk ingest request
            documentIndex = 0
            err := es.Bulk( conf.EsHost, conf.EsPort, conf.EsIndex, documents[0:] )
            if err != nil {
                fmt.Println( err )
                debug.PrintStack()
                continue
            }
        }
    }
    reader.Close()

    // send remaining documents
    err := es.Bulk( conf.EsHost, conf.EsPort, conf.EsIndex, documents[0:documentIndex] )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
    }
}

func ingestFiles( ingestInfo *IngestInfo ) {
    fmt.Printf( "Indexing %d files containing %d bytes of data\n",
        len( ingestInfo.Paths ),
        ingestInfo.NumTotalBytes )
    numIngestedBytes := int64( 0 )
    for i, path := range ingestInfo.Paths {
        // print current progress
        progress.PrintProgress( ingestInfo.NumTotalBytes, numIngestedBytes )
        // ingest the file
        ingestFile( path )
        // increment ingested bytes
        numIngestedBytes += ingestInfo.FileSizes[i]
    }
    // finished succesfully
    progress.PrintTotalDuration()
}

func parseArgs() Conf {
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
    return Conf{
        EsHost: *esHost,
        EsPort: *esPort,
        EsIndex: *esIndex,
        EsEndpoint: *esHost + ":" + *esPort + "/" + *esIndex,
        EsClearExisting: *esClearExisting,
        HdfsHost: *hdfsHost,
        HdfsPort: *hdfsPort,
        HdfsEndpoint: *hdfsHost + ":" + *hdfsPort,
        HdfsPath: *hdfsPath,
    }
}

type Conf struct {
    EsHost string
    EsPort string
    EsIndex string
    EsEndpoint string
    EsClearExisting bool
    HdfsHost string
    HdfsPort string
    HdfsEndpoint string
    HdfsPath string
}

var conf Conf

func main() {


    // parse flags
    flag.Parse()

    // load args into config struct
    conf = parseArgs()

    // check if index exists
    indexExists, err := es.IndexExists( conf.EsHost, conf.EsPort, conf.EsIndex )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }

    // if index exists
    if indexExists && conf.EsClearExisting {
        err = es.ClearIndex( conf.EsHost, conf.EsPort, conf.EsIndex )
        if err != nil {
            fmt.Println("Failed to delete index, aborting.")
            debug.PrintStack()
            os.Exit(1)
        }
    }

    // if index does not exist at this point
    if !indexExists || conf.EsClearExisting {
        err = es.CreateIndex( conf.EsHost, conf.EsPort, conf.EsIndex, `{
            "mappings": {
                "datum": {
                    "properties": {
                        "locality": {
                            "type": "object",
                            "properties": {
                                "location": {
                                    "type": "geo_point"
                                }
                            }
                        }
                    }
                }
            }
        }`)
        if err != nil {
            fmt.Println( err )
            debug.PrintStack()
            os.Exit(1)
        }
    }

    // get ingest info
    ingestInfo := getFileInfo()

    // ingest files
    ingestFiles( ingestInfo )
}
