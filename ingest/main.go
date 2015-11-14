package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "time"
    "net/http"

    "encoding/json"
    "runtime/debug"

    "github.com/parnurzeal/gorequest"
    "github.com/colinmarc/hdfs"
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

type Conf struct {
    esHost string
    esPort string
    esIndex string
    esEndpoint string
    esClearExisting bool
    hdfsHost string
    hdfsPort string
    hdfsEndpoint string
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
    percentComplete := 100 * ( float64( bytes ) / float64( totalBytes ) )
    bytesPerSecond := 1.0
    if elapsed > 0 {
        bytesPerSecond = float64( bytes ) / float64( elapsed )
    }
    estimatedSecondsRemaining := ( float64( totalBytes ) - float64( bytes ) ) / bytesPerSecond
    formattedTime := formatTime( uint64( estimatedSecondsRemaining ) )
    fmt.Printf( "\rIndexed %d bytes at %f Bps, %f%% complete, estimated time remaining %d:%02d:%02d",
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
func getHdfsClient() ( *hdfs.Client, error ) {
    if hdfsClient == nil {
        fmt.Println( "Connecting to HDFS: " + conf.hdfsEndpoint )
        client, err := hdfs.New( conf.hdfsEndpoint )
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
    client, err := getHdfsClient()
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }
    fmt.Println( "Retreiving ingest directory information from: " + conf.hdfsPath )
    files, err := client.ReadDir( conf.hdfsPath )
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
    fmt.Printf( "Indexing %d files containing %d bytes of data\n",
        len( ingestInfo.Paths ),
        ingestInfo.NumTotalBytes )
    startTime := uint64( time.Now().Unix() )
    numIngestedBytes := int64( 0 )
    for i, path := range ingestInfo.Paths {
        // print current progress
        printProgress( ingestInfo.NumTotalBytes, numIngestedBytes, startTime )
        // ingest the file
        ingestFile( path )
        // increment ingested bytes
        numIngestedBytes += ingestInfo.FileSizes[i]
    }
    // finished succesfully
    formattedTime := formatTime( uint64( time.Now().Unix() ) - startTime )
    fmt.Printf( "\nIndexing completed in %d:%02d:%02d\n",
        formattedTime.Hours,
        formattedTime.Minutes,
        formattedTime.Seconds )
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
    Userid string `json:"userid"`
    Username string `json:"username"`
    Hashtags []string `json:"hashtags"`
}

type TweetLocality struct {
    Timestamp string `json:"timestamp"`
    Label *string `json:"label"`
    Location *string `json:"location"`
}

type TweetSource struct {
    Properties TweetProperties `json:"properties"`
    Locality TweetLocality `json:"locality"`
}

type TweetIndex struct {
    Index string `json:"_index"`
    Type string `json:"_type"`
    Id string `json:"_id"`
}

type TweetIndexAction struct {
    Index *TweetIndex `json:"index"`
}

type IndexDocumentAction struct {
    IndexAction *TweetIndexAction
    Source *TweetSource
}

func columnExists( col string ) bool {
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
func buildIndexAction( tweetCsv []string ) *IndexDocumentAction {
    isoDate := tweetDateToISO( tweetCsv[0] )
    locality := TweetLocality {
        Timestamp: isoDate,
    }
    // state / province label may not exist
    if columnExists( tweetCsv[9] ) {
        locality.Label = &tweetCsv[9]
    }
    // long / lat may not exist
    if columnExists( tweetCsv[6] ) && columnExists( tweetCsv[7] ){
        location := ( tweetCsv[7] + "," + tweetCsv[6] )
        locality.Location = &location
    }
    properties := TweetProperties {
        Userid: tweetCsv[1],
        Username: tweetCsv[2],
        Hashtags: make([]string, 0),
    }
    // hashtags may not exist
    if ( columnExists( tweetCsv[5] ) ) {
        properties.Hashtags = strings.Split( strings.TrimSpace( tweetCsv[5] ), "#" )
    }
    // build source node
    source := TweetSource{
        Properties: properties,
        Locality: locality,
    }
    // build index
    index := TweetIndex{
        Index: conf.esIndex,
        Type: "datum",
        Id: tweetCsv[3],
    }
    // create index action
    indexAction := TweetIndexAction{
        Index: &index,
    }
    // return create bulk action
    return &IndexDocumentAction{
        IndexAction: &indexAction,
        Source: &source,
    }
}

func bulkUpload( documents []string ) {
    jsonLines := fmt.Sprintf( "%s\n", strings.Join( documents, "\n" ) )
	response, err := http.Post( conf.esHost + ":" + conf.esPort + "/_bulk", "application/json", strings.NewReader( jsonLines ) )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }
	response.Body.Close()
}

func ingestFile( filePath string ) {
    const batchSize = 9999
    documents := make( []string, batchSize )
    documentIndex := 0

    // get hdfs client
    client, err := getHdfsClient()
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }

    reader, err := client.Open( filePath )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }

    scanner := bufio.NewScanner( reader )
    for scanner.Scan() {
        line := strings.Split( scanner.Text(), "\t" )
        action := buildIndexAction( line )
        index, err := json.Marshal( action.IndexAction )
        if err != nil {
            fmt.Println( err )
            debug.PrintStack()
            continue
        }
        source, err := json.Marshal( action.Source )
        documents[ documentIndex ] = string( index ) + "\n" + string( source )
        documentIndex++

        if documentIndex % batchSize == 0 {
            // send bulk ingest request
            documentIndex = 0
            bulkUpload( documents[0:] )
        }
    }
    reader.Close()

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
    if *hdfsPort == "" {
        fmt.Println("HDFS port is not specified, please provide CL arg '-hdfs-port'.")
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
        hdfsPort: *hdfsPort,
        hdfsEndpoint: *hdfsHost + ":" + *hdfsPort,
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
        debug.PrintStack()
        os.Exit(1)
    }
    indexExists := resp.StatusCode != 404

    // if index exists
    if indexExists && conf.esClearExisting {
        fmt.Println( "Clearing index '" + conf.esIndex + "'." )
        _, _, errs := request.
            Delete( conf.esEndpoint ).
            End()
        if errs != nil {
            fmt.Println("Failed to delete index, aborting.")
            debug.PrintStack()
            os.Exit(1)
        }
    }

    if !indexExists || conf.esClearExisting {
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
            debug.PrintStack()
            os.Exit(1)
        }
    }

    ingestInfo := getFileInfo()
    ingestFiles( ingestInfo )
}
