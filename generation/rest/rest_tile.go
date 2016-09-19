package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultExt       = "json"
	defaultScheme    = "http"
	defaultBaseURL   = ""
	defaultIgnoreErr = false
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
	// get endpoint
	endpoint, ok := json.GetString(g.req.Params, "endpoint")
	if !ok {
		return nil, fmt.Errorf("Missing `endpoint` parameter")
	}
	// whether to ingore error of not
	ignoreErr := json.GetBoolDefault(g.req.Params, defaultIgnoreErr, "ignoreErr")
	// get scheme
	scheme := json.GetStringDefault(g.req.Params, defaultScheme, "scheme")
	// get ext
	ext := json.GetStringDefault(g.req.Params, defaultExt, "ext")
	// create URL
	url := fmt.Sprintf("%s://%s/%s/%d/%d/%d.%s",
		scheme,
		endpoint,
		g.req.URI,
		g.req.Coord.Z,
		g.req.Coord.X,
		g.req.Coord.Y,
		ext)
	fmt.Println(url)
	// build http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// set appropriate headers based on extention
	if ext == "bin" {
		req.Header.Set("Accept", "application/octet-stream")
	} else {
		req.Header.Set("Accept", "application/json")
	}
	// build http request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// read result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// check status code
	if res.StatusCode >= 300 {
		if ignoreErr {
			return []byte{}, nil
		}
		return nil, fmt.Errorf(string(body))
	}
	return body, nil
}
