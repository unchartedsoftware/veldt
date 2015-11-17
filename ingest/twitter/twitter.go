package twitter

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
    "time"
    "encoding/json"
    "runtime/debug"

    "github.com/unchartedsoftware/prism/binning"

    "github.com/unchartedsoftware/prism/ingest/conf"
    "github.com/unchartedsoftware/prism/ingest/es"
    "github.com/unchartedsoftware/prism/ingest/hdfs"
)

type TweetProperties struct {
    Userid string `json:"userid"`
    Username string `json:"username"`
    Hashtags []string `json:"hashtags"`
}

type TweetLocality struct {
    Timestamp string `json:"timestamp"`
    Location *binning.LonLat `json:"location"`
    Pixel *binning.PixelCoord `json:"pixel"`
}

type TweetSource struct {
    Properties TweetProperties `json:"properties"`
    Locality TweetLocality `json:"locality"`
    Text string `json:"text"`
}

type TweetIndex struct {
    Type string `json:"_type"`
    Id string `json:"_id"`
}

type TweetIndexAction struct {
    Index *TweetIndex `json:"index"`
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

func columnExists( col string ) bool {
    if col != "" && col != "None" {
        return true
    }
    return false
}

const maxLevelSupported = 24
const tileResolution = 256

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
func createIndexAction( tweetCsv []string ) ( *string, error ) {
    isoDate := tweetDateToISO( tweetCsv[0] )
    locality := TweetLocality {
        Timestamp: isoDate,
    }
    // long / lat may not exist
    if columnExists( tweetCsv[6] ) && columnExists( tweetCsv[7] ) {
        lon, lonErr := strconv.ParseFloat( tweetCsv[6], 64 )
        lat, latErr := strconv.ParseFloat( tweetCsv[7], 64 )
        if lonErr == nil && latErr == nil {
            lonLat := &binning.LonLat{
                Lat: lat,
                Lon: lon,
            }
            pixel := binning.LonLatToPixelCoord( lonLat, maxLevelSupported, tileResolution );
            locality.Location = lonLat
            locality.Pixel = pixel
        }
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
        Text: tweetCsv[4],
    }
    // build index
    index := TweetIndex{
        Type: "datum",
        Id: tweetCsv[3],
    }
    // create index action
    indexAction := TweetIndexAction{
        Index: &index,
    }
    indexBytes, indexErr := json.Marshal( indexAction )
    if indexErr != nil {
        return nil, indexErr
    }
    sourceBytes, sourceErr := json.Marshal( source )
    if sourceErr != nil {
        return nil, sourceErr
    }
    jsonString := string( indexBytes ) + "\n" + string( sourceBytes )
    return &jsonString, nil
}

const Mappings =
    `{
        "datum": {
            "properties": {
                "locality": {
                    "type": "object",
                    "properties": {
                        "location": {
                            "type": "geo_point"
                        },
                        "hashtags" : {
                          "type" : "string",
                          "index" : "not_analyzed"
                        },
                        "userid" : {
                          "type" : "string",
                          "index" : "not_analyzed"
                        },
                        "username" : {
                          "type" : "string",
                          "index" : "not_analyzed"
                        }
                    }
                }
            }
        }
    }`

func TwitterWorker( file os.FileInfo ) {
    config := conf.GetConf()
    documents := make( []string, config.BatchSize )
    documentIndex := 0

    // get hdfs client
    client, clientErr := hdfs.GetHdfsClient( config.HdfsHost, config.HdfsPort )
    if clientErr != nil {
        fmt.Println( clientErr )
        debug.PrintStack()
        return
    }

    // get hdfs file reader
    reader, fileErr := client.Open( config.HdfsPath + "/" + file.Name() )
    if fileErr != nil {
        fmt.Println( fileErr )
        debug.PrintStack()
        return
    }

    scanner := bufio.NewScanner( reader )
    for scanner.Scan() {
        line := strings.Split( scanner.Text(), "\t" )
        action, err := createIndexAction( line )
        if err != nil {
            fmt.Println( err )
            debug.PrintStack()
            continue
        }
        documents[ documentIndex ] = *action
        documentIndex++
        if documentIndex % config.BatchSize == 0 {
            // send bulk ingest request
            documentIndex = 0
            err := es.Bulk( config.EsHost, config.EsPort, config.EsIndex, documents[0:] )
            if err != nil {
                fmt.Println( err )
                debug.PrintStack()
                continue
            }
        }
    }
    reader.Close()

    // send remaining documents
    err := es.Bulk( config.EsHost, config.EsPort, config.EsIndex, documents[0:documentIndex] )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
    }
}
