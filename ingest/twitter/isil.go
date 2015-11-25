package twitter

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/unchartedsoftware/prism/binning"
)

// ISILTweetDocument represents a single TSV row of isil keyword twitter data.
type ISILTweetDocument struct {
	Cols []string
}

// SetData sets the internal TSV column.
func (d *ISILTweetDocument) SetData(cols []string) {
	d.Cols = cols
}

// GetID returns the document id.
func (d ISILTweetDocument) GetID() string {
	return d.Cols[0]
}

// GetType returns the document type.
func (d ISILTweetDocument) GetType() string {
	return "datum"
}

// GetMappings returns the documents mappings.
func (d ISILTweetDocument) GetMappings() string {
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

// ISILSource is the source structure for this document.
type ISILSource struct {
	UserID    *string             `json:"userid"`
	Username  string              `json:"username"`
	Hashtags  []string            `json:"hashtags"`
	URLs      []string            `json:"hashtags"`
	Timestamp string              `json:"timestamp"`
	Text      string              `json:"text"`
	LonLat    *binning.LonLat     `json:"lonlat"`
	Pixel     *binning.PixelCoord `json:"pixel"`
	Rankings  map[string]uint64   `json:"rankings"`
}

// GetSource returns the marshalled source portion of the document.
func (d ISILTweetDocument) GetSource() ([]byte, error) {
	// CSV line as array:
	//	  0: tweet id
	//	  1: tweet datetime
	//	  2: account name
	//	  3: user-provided real name
	//	  4: latitude
	//	  5: longitude
	//	  6: text of tweet
	//	  7: language
	//	  8: source of tweet
	//	  9: #-delimited list of hashtags
	//	  10: this tweet has been favorited
	//	  11: number of times this tweet has been favorited
	//	  12: has been retweeted
	//	  13: number of times this tweet has been retweeted
	//	  14: number of times this retweet has been favorited
	//	  15: number of times this retweet has been retweeted
	//	  16: unsorted (user-provided) url in text
	//	  17: user-provided country
	//	  18: user-provided place type
	//	  19: @screen_name mentions in text
	//	  20: screen_name this tweet is a reply to
	//	  21: tweet id this tweets is a reply to
	//	  22: comma-delimited media metadata type + media url
	//	  23: numeric id of user's account
	//	  24: number of tweets user has made
	//	  25: number of followers user has
	//	  26: number of friends user has (number of accounts this user follows)
	//	  27: user-provided description of themself
	//	  28: datetime of account creation
	//	  29: is geo-location enabled
	//	  30: number of lists user is a member of
	//	  31: user-provided location
	//	  32: user-provided url of banner image
	//	  33: user-provided time zone information
	//	  34: user-provided personal url
	//	  35: has account been verified by Twitter
	//	  36: the 'real' name given by the user
	cols := d.Cols
	timestamp, err := tweetDateToISO(cols[1])
	if err != nil {
		return nil, err
	}
	source := &ISILSource{
		Username:  cols[2],
		Hashtags:  make([]string, 0),
		URLs:      make([]string, 0),
		Timestamp: timestamp,
		Text:      cols[6],
		Rankings:  make(map[string]uint64),
	}
	// user id may not exist
	if columnExists(cols[23]) {
		source.UserID = &cols[23]
	}
	// lon / lat data may not exist
	if columnExists(cols[4]) && columnExists(cols[5]) {
		lon, lonErr := strconv.ParseFloat(cols[5], 64)
		lat, latErr := strconv.ParseFloat(cols[4], 64)
		if lonErr == nil && latErr == nil {
			source.LonLat = &binning.LonLat{
				Lat: lat,
				Lon: lon,
			}
			source.Pixel = binning.LonLatToPixelCoord(source.LonLat)
		}
	}
	// hashtags may not exist
	if columnExists(cols[9]) {
		source.Hashtags = strings.Split(strings.TrimSpace(cols[9]), "#")
	}
	// URLs may not exist
	if columnExists(cols[16]) {
		source.URLs = strings.Split(strings.TrimSpace(cols[16]), ",")
	}
	return json.Marshal(source)
}
