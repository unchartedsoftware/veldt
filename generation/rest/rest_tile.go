package rest

import (
	"io/ioutil"
	"net/http"
	"errors"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
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
	// build http request
	client := &http.Client{}
	url, exists := json.Get(g.req.Params, "url");
	if exists == false {
		return nil, errors.New("Url missing from request params")
	}
	req, err := http.NewRequest("GET", url.(string), nil)
	if err != nil {
		return nil, err
	}
	// set appropriate headers
	var contentType string
	if json.GetStringDefault(g.req.Params, "json", "extension") == "bin" {
		contentType = "application/octet-stream"
	} else {
		contentType = "application/json"
	}
	req.Header.Set("Accept", contentType)
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
	return body, nil
}
