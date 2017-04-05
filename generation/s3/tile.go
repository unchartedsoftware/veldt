package s3

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/util/json"
)

const (
	defaultExt       = "bin"
	defaultScheme    = "http"
	defaultBaseURL   = ""
	defaultPadCoords = true
	defaultKeyPrefix = ""
)

// Tile represents an S3 tile type
type Tile struct {
	padCoords bool
	ext       string
}

// NewTile instantiates and returns a new S3 tile.
func NewTile() veldt.TileCtor {
	return func() (veldt.Tile, error) {
		return &Tile{}, nil
	}
}

// Parse parses the provided JSON object for tile attributes
func (t *Tile) Parse(params map[string]interface{}) error {
	// Get file extension
	t.ext = json.GetStringDefault(params, defaultExt, "ext")
	// Whether of not to padd the coordinates with leading 0s
	t.padCoords = json.GetBoolDefault(params, defaultPadCoords, "padCoords")
	return nil
}

// Create generates a tile from the provided URI, tile coordinate and query parameters.
func (t *Tile) Create(s3uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// create s3 client
	s3Client, err := NewS3Client()
	if err != nil {
		return nil, err
	}
	// Get s3 bucket & key. g.req.URI is of the form: "{bucket-name}/{key-path/maybe/containing/slashes}"
	uri := strings.SplitN(s3uri, "/", 2)
	bucket := uri[0]
	prefix := uri[1]

	// Format key (s3 filename)
	format := "%s/%d/%d/%d.%s"
	if t.padCoords {
		digits := strconv.Itoa(int(math.Floor(math.Log10(float64(int(1)<<coord.Z)))) + 1)
		format = "%s/%02d/%0" + digits + "d/%0" + digits + "d.%s"
	}
	// Construct the key
	key := fmt.Sprintf(format,
		prefix,
		coord.Z,
		coord.X,
		coord.Y,
		t.ext)

	// Create request
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	// Fetch tile from s3
	res, err := s3Client.GetObject(params)
	// Handle response
	if err != nil {
		// don't return an error if the tile doesn't exist
		if err.Error() == s3.ErrCodeNoSuchKey {
			return []byte{}, nil
		}
		return nil, err
	}
	// read result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
