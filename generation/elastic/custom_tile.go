package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/tile"
	jsonp "github.com/unchartedsoftware/prism/util/json"
)

// CustomTile represents a tiling generator that produces top term counts.
type CustomTile struct {
	TileGenerator
	Tiling     *param.Tiling
	Source     map[string]interface{}
}

// NewCustomTile instantiates and returns a pointer to a new generator.
func NewCustomTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		elastic, err := param.NewElastic(tileReq)
		if err != nil {
			return nil, err
		}
		// required
		tiling, err := param.NewTiling(tileReq)
		if err != nil {
			return nil, err
		}
		source, ok := jsonp.GetChild(tileReq.Params, "source")
		if !ok {
			return nil, fmt.Errorf("Source was not of type `map[string]interface{}` in response for request %s",
				tileReq.String())
		}
		t := &CustomTile{}
		t.Elastic = elastic
		t.Tiling = tiling
		t.Source = source
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *CustomTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
	}
}

func (g *CustomTile) GetQuerySource() (interface{}, error) {
	xFilter, err := g.Tiling.GetXQuery().Source()
	if err != nil {
		return nil, err
	}
	yFilter, err := g.Tiling.GetYQuery().Source()
	if err != nil {
		return nil, err
	}
	source := g.Source
	must, ok := jsonp.GetArray(source, "must")
	if !ok {
		must = make([]interface{}, 0)
	}
	must = append(must, xFilter)
	must = append(must, yFilter)
	source["must"] = must
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": source,
		},
	}
	return query, nil
}

func (g *CustomTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	// Return the raw results
	return json.Marshal(res)
}

// GetTile returns the marshalled tile data.
func (g *CustomTile) GetTile() ([]byte, error) {
	source, err := g.GetQuerySource()
	if err != nil {
		return nil, err
	}
	// send query
	res, err := g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Size(0).
		Source(source).
		Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
