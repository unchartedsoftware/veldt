package s3

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	s3BaseURL = "https://s3.amazonaws.com/xdata-tiles"
)

// TopicCountTile represents a tiling generator that produces heatmaps.
type TopicCountTile struct {
	TileGenerator
}

// NewTopicCountTile instantiates and returns a pointer to a new generator.
func NewTopicCountTile() tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		t := &TopicCountTile{}
		t.baseURL = s3BaseURL
		t.req = tileReq
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopicCountTile) GetParams() []tile.Param {
	return []tile.Param{}
}

// GetTile returns the marshalled tile data.
func (g *TopicCountTile) GetTile() ([]byte, error) {
	// send query
	client := &http.Client{}

	/* g.req.Index(s3 bucket name) is passed in from client.
	Typically it has a slash (ex: census-hackathon-2016/types-word-cloud)
	however having this as part of the request means the prism-server
	interprets it as two parameters. Replaced '/' with ':' on the client side,
	and the following reverses this action.
	*/
	bucketName := strings.Replace(g.req.Index, ":", "/", 1)
	url := g.baseURL + "/" + bucketName + "/" + strconv.Itoa(int(g.req.Coord.X)) + "/" + strconv.Itoa(int(g.req.Coord.Y)) + "/" + strconv.Itoa(int(g.req.Coord.Z)) + ".json"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// parse and return results
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return body, nil
}
