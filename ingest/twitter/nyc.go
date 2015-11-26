package twitter

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/unchartedsoftware/prism/binning"
)

// NYCTweetDocument represents a single TSV row of nyc twitter data.
type NYCTweetDocument struct {
	Cols []string
}

// Setup initialized and required state prior to ingestion.
func (d NYCTweetDocument) Setup() error {
	return nil
}

// Teardown cleans up any state after ingestion.
func (d NYCTweetDocument) Teardown() error {
	return nil
}

// SetData sets the internal TSV column.
func (d *NYCTweetDocument) SetData(cols []string) {
	d.Cols = cols
}

// GetID returns the document id.
func (d NYCTweetDocument) GetID() string {
	return d.Cols[3]
}

// GetType returns the document type.
func (d NYCTweetDocument) GetType() string {
	return "datum"
}

// GetMappings returns the documents mappings.
func (d NYCTweetDocument) GetMappings() string {
	return `{
        "` + d.GetType() + `": {
            "location": {
                "type": "geo_point"
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
    }`
}

// NYCSource is the source structure for this document.
type NYCSource struct {
	UserID    string              `json:"userid"`
	Username  string              `json:"username"`
	Hashtags  []string            `json:"hashtags"`
	Timestamp string              `json:"timestamp"`
	Text      string              `json:"text"`
	LonLat    *binning.LonLat     `json:"lonlat"`
	Pixel     *binning.PixelCoord `json:"pixel"`
}

// GetSource returns the marshalled source portion of the document.
func (d NYCTweetDocument) GetSource() ([]byte, error) {
	// CSV line as array:
	//     0: 'Fri Jan 04 18:42:42 +0000 2013',
	//     1: '242573761',
	//     2: 'AdioAsh5',
	//     3:  '287267829735100416',
	//     4:  "Blah blah blah blah blah",
	//     5:  '',
	//     6:  '-73.94068643', {lon}
	//     7:  '40.66179087', {lat}
	//     8:  'United States',
	//     9:  'New York, NY',
	//     10:  'city',
	//     11:  'en'
	cols := d.Cols
	timestamp, err := tweetDateToISO(cols[0])
	if err != nil {
		return nil, err
	}
	source := &NYCSource{
		UserID:    cols[1],
		Username:  cols[2],
		Hashtags:  make([]string, 0),
		Timestamp: timestamp,
		Text:      cols[4],
	}
	// lon / lat data may not exist
	if columnExists(cols[6]) && columnExists(cols[7]) {
		lon, lonErr := strconv.ParseFloat(cols[6], 64)
		lat, latErr := strconv.ParseFloat(cols[7], 64)
		if lonErr == nil && latErr == nil {
			source.LonLat = &binning.LonLat{
				Lat: lat,
				Lon: lon,
			}
			source.Pixel = binning.LonLatToPixelCoord(source.LonLat)
		}
	}
	// hashtags may not exist
	if columnExists(cols[5]) {
		source.Hashtags = strings.Split(strings.TrimSpace(cols[5]), "#")
	}
	return json.Marshal(source)
}
