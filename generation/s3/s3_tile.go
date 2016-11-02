package s3

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/unchartedsoftware/prism/tile"
	awsManager "github.com/unchartedsoftware/prism/util/aws"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultExt       = "json"
	defaultScheme    = "http"
	defaultBaseURL   = ""
	defaultIgnoreErr = false
	defaultPadCoords = false
	defaultKeyPrefix = ""
)

// Tile represents a tiling generator that produces heatmaps.
type Tile struct {
	TileGenerator
}

// NewTile instantiates and returns a pointer to a new generator.
func NewTile() tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		t := &Tile{}
		t.req = tileReq
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *Tile) GetParams() []tile.Param {
	return []tile.Param{}
}

// GetTile returns the marshalled tile data.
func (g *Tile) GetTile() ([]byte, error) {
	// create s3 client
	s3Client, err := awsManager.NewS3Client()
	if err != nil {
		return nil, err
	}
	// get path
	keyPrefix := json.GetStringDefault(g.req.Params, defaultKeyPrefix, "keyPrefix")
	// whether to ingore error or not
	ignoreErr := json.GetBoolDefault(g.req.Params, defaultIgnoreErr, "ignoreErr")
	// get ext
	ext := json.GetStringDefault(g.req.Params, defaultExt, "ext")
	// get padCoords
	padCoords := json.GetBoolDefault(g.req.Params, defaultPadCoords, "padCoords")
	// Format key (s3 filename)
	format := "%d/%d/%d.%s"
	if padCoords {
		digits := strconv.Itoa(int(math.Floor(math.Log10(float64(int(1)<<g.req.Coord.Z)))) + 1)
		format = "%02d/%0" + digits + "d/%0" + digits + "d.%s"
	}
	if keyPrefix != defaultKeyPrefix {
		format = "%s/" + format
	}
	key := fmt.Sprintf(format,
		keyPrefix,
		g.req.Coord.Z,
		g.req.Coord.X,
		g.req.Coord.Y,
		ext)
	// Create request
	params := &s3.GetObjectInput{
		Bucket: aws.String(g.req.URI),
		Key:    aws.String(key),
	}
	res, err := s3Client.GetObject(params)
	// Handle response
	if err != nil {
		if ignoreErr {
			return []byte{}, nil
		}
		return nil, fmt.Errorf(err.Error())
	}
	// read result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
