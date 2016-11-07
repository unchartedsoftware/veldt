package tile

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tile"
)

const (
	xAggName = "x"
	yAggName = "y"
)

// Heatmap represents a tiling generator that produces heatmaps.
type Heatmap struct {
	client *elastic.Client
}

// NewHeatmap instantiates a new heatmap tile.
func NewHeatmap(host, port string) tile.GeneratorConstructor {
	return func() (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		return &Heatmap{
			client: client,
		}, nil
	}
}

func (g *Heatmap) GetQuery(q query.Query) (elastic.Query, error) {
	root := elastic.NewBoolQuery()
	root.Must(q.GetXQuery())
	root.Must(q.GetYQuery())
	if q != nil {
		err := q.Apply(root)
		if err != nil {
			return nil, err
		}
	}
	return root, nil
}

func (g *Heatmap) GetTile(req *tile.Request) ([]byte, error) {
	// create search service
	search := g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Size(0).
		Query(g.GetQuery(req.Query))
	// apply agg
	g.ApplyXYAgg(search)
	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	bins := g.GetXYBins(res)
	buffer := make([]float64, len(bins))
	for i, bin := range bins {
		buffer[i] = bin.DocCount
	}
	return g.Float64ToBytes(buffer), nil
}
