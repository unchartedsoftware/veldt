package twitter

import (
    "fmt"
    "strings"
    "strconv"
    "time"
    "encoding/json"
)

type TweetProperties struct {
    Userid string `json:"userid"`
    Username string `json:"username"`
    Hashtags []string `json:"hashtags"`
}

type TweetLonLat struct {
    Lon float64 `json:"lat"`
    Lat float64 `json:"lon"`
}

type TweetXY struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
}

type TweetLocality struct {
    Timestamp string `json:"timestamp"`
    Label *string `json:"label"`
    Location *TweetLonLat `json:"location"`
    XY *TweetXY `json:"xy"`
}

type TweetSource struct {
    Properties TweetProperties `json:"properties"`
    Locality TweetLocality `json:"locality"`
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
func CreateIndexAction( tweetCsv []string ) ( *string, error ) {
    isoDate := tweetDateToISO( tweetCsv[0] )
    locality := TweetLocality {
        Timestamp: isoDate,
    }
    // state / province label may not exist
    if columnExists( tweetCsv[9] ) {
        locality.Label = &tweetCsv[9]
    }
    // long / lat may not exist
    if columnExists( tweetCsv[6] ) && columnExists( tweetCsv[7] ) {
        lon, lonErr := strconv.ParseFloat( tweetCsv[6], 64 )
        lat, latErr := strconv.ParseFloat( tweetCsv[7], 64 )
        if lonErr == nil && latErr == nil {
            location := &TweetLonLat{
                Lat: lat,
                Lon: lon,
            }
            xy := &TweetXY{
                X: 0,
                Y: 0,
            }
            locality.Location = location
            locality.XY = xy
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
