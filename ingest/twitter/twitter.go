package twitter

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "encoding/json"
    "runtime/debug"

    "github.com/unchartedsoftware/prism/binning"

    "github.com/unchartedsoftware/prism/ingest/conf"
    "github.com/unchartedsoftware/prism/ingest/es"
    "github.com/unchartedsoftware/prism/ingest/hdfs"
    "github.com/unchartedsoftware/prism/ingest/terms"
)

type TweetProperties struct {
    UserID string `json:"userid"`
    Username string `json:"username"`
    Hashtags []string `json:"hashtags"`
    TopTerms []string `json:"topterms"`
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
    ID string `json:"_id"`
}

type TweetIndexAction struct {
    Index *TweetIndex `json:"index"`
}

const maxLevelSupported = 24
const tileResolution = 256

func createIndexAction( tweet *TweetData ) ( *string, error ) {
    locality := TweetLocality {
        Timestamp: tweet.ISODate,
    }
    // long / lat may not exist
    if tweet.LonLat != nil {
        lonLat := tweet.LonLat
        pixel := binning.LonLatToPixelCoord( lonLat, maxLevelSupported, tileResolution );
        locality.Location = lonLat
        locality.Pixel = pixel
    }
    properties := TweetProperties {
        UserID: tweet.UserID,
        Username: tweet.Username,
        Hashtags: tweet.Hashtags,
        TopTerms: terms.GetTopTerms( tweet.TweetText ),
    }
    // build source node
    source := TweetSource{
        Properties: properties,
        Locality: locality,
        Text: tweet.TweetText,
    }
    // build index
    index := TweetIndex{
        ID: tweet.TweetID,
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

func TwitterWorker( file os.FileInfo ) {
    config := conf.GetConf()
    documents := make( []string, config.BatchSize )
    documentIndex := 0

    // get hdfs file reader
    reader, err := hdfs.Open( config.HdfsHost, config.HdfsPort, config.HdfsPath + "/" + file.Name() )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        return
    }

    scanner := bufio.NewScanner( reader )
    for scanner.Scan() {
        line := strings.Split( scanner.Text(), "\t" )
        tweet := ParseTweetData( line )
        action, err := createIndexAction( tweet )
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
            err := es.Bulk( config.EsHost, config.EsPort, config.EsIndex, config.EsType, documents[0:] )
            if err != nil {
                fmt.Println( err )
                debug.PrintStack()
                continue
            }
        }
    }
    reader.Close()

    // send remaining documents
    err = es.Bulk( config.EsHost, config.EsPort, config.EsIndex, config.EsType, documents[0:documentIndex] )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
    }
}

func TwitterTopTermsWorker( file os.FileInfo ) {
    config := conf.GetConf()

    // get hdfs file reader
    reader, err := hdfs.Open( config.HdfsHost, config.HdfsPort, config.HdfsPath + "/" + file.Name() )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        return
    }

    scanner := bufio.NewScanner( reader )
    for scanner.Scan() {
        line := strings.Split( scanner.Text(), "\t" )
        tweet := ParseTweetData( line )
        terms.AddTerms( tweet.TweetText )
    }

    reader.Close()
}
