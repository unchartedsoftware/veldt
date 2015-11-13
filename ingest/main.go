package main

import (
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/parnurzeal/gorequest"
    "github.com/colinmarc/hdfs"
)

var (
	esHost = flag.CommandLine.String( "es-host", "", "Elasticsearch host" )
    esPort = flag.CommandLine.String( "es-port", "9200", "Elasticsearch port" )
	esIndex = flag.CommandLine.String( "es-index", "", "Elasticsearch index" )
	esClearExisting = flag.CommandLine.Bool( "es-clear-existing", true, "Clear index before ingest" )

    hdfsHost = flag.CommandLine.String( "hdfs-host", "", "HDFS host" )
    hdfsPath = flag.CommandLine.String( "hdfs-path", "", "HDFS ingest source data path" )
)

type Conf struct {
    esHost string
    esPort string
    esIndex string
    esEndpoint string
    esClearExisting bool
    hdfsHost string
    hdfsPath string
}

type Time struct {
    Seconds uint64
    Minutes uint64
    Hours uint64
}

func formatTime( totalSeconds uint64 ) Time {
    totalMinutes := totalSeconds / 60
    seconds := totalSeconds % 60
    hours := totalMinutes / 60
    minutes := totalMinutes % 60
    return Time{
        Seconds: seconds,
        Minutes: minutes,
        Hours: hours,
    }
}

func printProgress( totalBytes int64, bytes int64, startTime uint64 ) {
    elapsed := uint64( time.Now().Unix() ) - startTime
    percentComplete := 100 * ( bytes / totalBytes )
    bytesPerSecond := float64( bytes ) / float64( elapsed )
    if bytesPerSecond == 0 {
        bytesPerSecond = 1
    }
    estimatedSecondsRemaining := ( float64( totalBytes ) - float64( bytes ) ) / bytesPerSecond
    formattedTime := formatTime( uint64( estimatedSecondsRemaining ) )
    fmt.Printf( "\rIndexed %d bytes at %d Bps, %d%% complete, estimated time remaining %d:%02d:%02d",
        bytes,
        bytesPerSecond,
        percentComplete,
        formattedTime.Hours,
        formattedTime.Minutes,
        formattedTime.Seconds )
}

func printTimeout( duration uint32 ) {
    for duration >= 0 {
        fmt.Printf( "\rRetrying in " + string( duration ) + " seconds..." )
        time.Sleep( time.Second )
        duration -= 1
    }
    fmt.Println()
}

var hdfsClient *hdfs.Client = nil
func getHdfsClient( host string ) ( *hdfs.Client, error ) {
    if hdfsClient == nil {
        fmt.Println( "Connecting to HDFS " + host )
        client, err := hdfs.New( conf.hdfsHost )
        hdfsClient = client
        return hdfsClient, err
    }
    return hdfsClient, nil
}

type IngestInfo struct {
    Paths []string
    FileSizes []int64
    NumTotalBytes int64
}

func getFileInfo() *IngestInfo {
    client, err := getHdfsClient( conf.hdfsHost )
    if err != nil {
        fmt.Println( err )
        os.Exit(1)
    }
    print( "Retreiving file information for " + conf.hdfsPath )
    files, err := client.ReadDir( conf.hdfsPath )
    if err != nil {
        fmt.Println( err )
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
            paths = append( paths, conf.hdfsPath + "/" + file.Name() )
            fileSizes = append( fileSizes, file.Size() )
        }
    }
    return &IngestInfo{
        Paths: paths,
        FileSizes: fileSizes,
        NumTotalBytes: numTotalBytes,
    }
}

func ingestFiles( ingestInfo *IngestInfo ) {
    fmt.Printf( "Indexing %d files containing %d bytes of data\n" )
    startTime := uint64( time.Now().Unix() )
    // endless loop of hell
    for {
        numIngestedBytes := 0
        for i, path := range ingestInfo.Paths {
            // print current progress
            printProgress( ingestInfo.NumTotalBytes, numIngestedBytes, startTime )
            // ingest the file
            err := ingestFile( path )
            if err != nil {
                fmt.Println( err )
                printTimeout( 5 )
                continue
            }
            // increment ingested bytes
            numIngestedBytes += ingestInfo.FileSizes[i]
        }
        // finished succesfully
        formattedTime = formatTime( uint64( time.Now().Unix() ) - startTime )
        fmt.Println( "\nIndexing completed in %d:%02d:%02d",
            formattedTime.Hours,
            formattedTime.Minutes,
            formattedTime.Seconds )
        break
    }
}

func tweetDateToISO( tweetDate string ) string {
    const layout = "Mon Jan 2 15:04:05 -0700 2006"
    t, err := time.Parse( layout, tweetDate )
    if err != nil {
        fmt.Println( "Error parsing date: " + tweetDate )
        return ""
    }
    return t.Format( time.RFC3339 )
}

type TweetProperties struct {
    Userid string `json:userid`
    Username string `json:username`
    Hashtags []string `json:hashtags`
}

type TweetLocality struct {
    Timestamp string `json:timestamp`
    Label *string `json:label`
    Location *string `json:location`
}

type TweetSource struct {
    properties TweetProperties `json:properties`
    locality TweetLocality `json:locality`
}

type TweetDocument struct {
    Index string `json:_index`
    Type string `json:_type`
    Id string `json:_id`
    Source TweetSource `json:_source`
}

type CreateDocument struct {
    Create TweetDocument `json:create`
}

func columnExists( string col ) bool {
    if col != "" && col != "None" {
        return true
    }
    return false
}

/*
    CSV line as array:
        0: 'Fri Jan 04 18:42:42 +0000 2013',
        1: '242573761',
        2: 'AdioAsh5',
        3:  '287267829735100416',
        4:  "Blah blah blah blah blah",
        5:  '',
        6:  '-73.94068643', {lon}
        7:  '40.66179087', {lat}
        8:  'United States',
        9:  'New York, NY',
        10:  'city',
        11:  'en'
*/
func buildTweetDocument( tweetCsv ) {
    isoDate := tweetDateToISO( tweetCsv[0] )
    if err != nil {
        fmt.Println(err)
    }
    locality := TweetLocality {
        Timestamp: isoDate,
    }
    // state / province label may not exist
    if columnExists( tweetCsv[9] ) {
        locality.Label = &tweetCsv[9]
    }
    // long / lat may not exist
    if columnExists( tweetCsv[6] ) && columnExists( tweetCsv[7] ){
        locality.Location = &( tweetCsv[7] + "," + tweetCsv[6] )
    }
    properties := TweetProperties {
        Userid: tweetCsv[1],
        Username: tweetCsv[2],
        Hashtags: strings.Split( strings.Trim( tweet[5] ), "#" ),
    }
    // build source node
    source := TweetSource{
        properties: properties,
        locality: locality,
    }
    // build document
    document := TweetDocument{
        Index: conf.esIndex,
        Type: "datum",
        Id: tweetCsv[3],
        Source: source,
    }
    // return create bulk action
    return CreateDocument {
        Create: TweetDocument,
    }
}

func bulkUpload( documents []string ) {
    jsonLines := strings.join( documents, '\n' )
    _, _, errs := request.
        Post( conf.esHost + ":" + conf.esHost + "/_bulk" ).
        Send( jsonLines ).
        End()
    if errs != nil {
        fmt.Println( errs )
        os.Exit(1)
    }
}

func ingestFile( filePath ) {
    const batchSize = 9999
    documents := make( []string, batchSize )
    documentIndex := 0

    // endless loop of death
    for {
        // get hdfs client
        client := getHdfsClient( conf.hdfsHost )
        if err != nil {
            fmt.Println( err )
            printTimeout( 5 )
            continue
        }

        reader, err := client.Open( filePath )
        if err != nil {
            fmt.Println( err )
            printTimeout( 5 )
            continue
        }

        scanner := bufio.NewScanner( reader )
        for scanner.Scan() {
            line := strings.split( scanner.Text(), '\t' )
            document := buildTweetDocument( line )
            documents[ documentIndex ] = json.Marshal( document )
            documentIndex++
        }

        if i % batchSize == 0 {
            // send bulk ingest request
            documentIndex = 0
            bulkUpload( documents[0:] )
        }
    }

    // send remaining documents
    bulkUpload( documents[0:documentIndex] )
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
    if *hdfsPath == "" {
        fmt.Println("HDFS path is not specified, please provide CL arg '-hdfs-path'.")
        os.Exit(1)
    }
    return Conf{
        esHost: *esHost,
        esPort: *esPort,
        esIndex: *esIndex,
        esEndpoint: *esHost + ":" + *esPort + "/" + *esIndex,
        esClearExisting: *esClearExisting,
        hdfsHost: *hdfsHost,
        hdfsPath: *hdfsPath,
    }
}

var conf Conf

func main() {

    // parse flags
    flag.Parse()
    // load args into config struct
    conf = parseArgs()

    // check if index exists
    request := gorequest.New()
    resp, _, errs := request.
		Head( conf.esEndpoint ).
		End()
    if errs != nil {
        fmt.Println( errs )
        os.Exit(1)
    }
    indexExists := resp.StatusCode != 404

    // if index exists
    if indexExists && esClearExisting {
        fmt.Println( "Clearing index '" + conf.esIndex + "'." )
        _, errs := request.
            Delete( conf.esEndpoint ).
            End()
        if errs != nil {
            fmt.Println("Failed to delete index, aborting.")
            os.Exit()
        }
    }

    if !indexExists {
        fmt.Println( "Creating index '" + conf.esIndex + "'." )
        mappingsBody := `{
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
        }`
        _, _, errs := request.
    		Put( conf.esEndpoint ).
            Send( mappingsBody ).
    		End()
        if errs != nil {
            fmt.Println( errs )
            os.Exit(1)
        }
    }

    ingestInfo := getFileInfo()
    ingestFiles( ingestInfo )
    fmt.Println()
}
