package elastic

import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/util/json"
)

// BinnedTopHits represents an elasticsearch implementation of the binned top
// hits tile.
type BinnedTopHits struct {
	Elastic
	Bivariate
	TopHits
}

// NewBinnedTopHits instantiates and returns a new tile struct.
func NewBinnedTopHits(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		b := &BinnedTopHits{}
		b.Host = host
		b.Port = port
		return b, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (b *BinnedTopHits) Parse(params map[string]interface{}) error {
	err := b.TopHits.Parse(params)
	if err != nil {
		return err
	}
	return b.Bivariate.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (b *BinnedTopHits) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// create search service
	search, err := b.CreateSearchService(uri)
	if err != nil {
		return nil, err
	}

	// create root query
	q, err := b.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q.Must(b.Bivariate.GetQuery(coord))
	// set the query
	search.Query(q)

	// get aggs
	topHitsAggs := b.TopHits.GetAggs()
	aggs := b.Bivariate.GetAggsWithNested(coord, "top-hits", topHitsAggs["top-hits"])

	// set the aggregation
	search.Aggregation("x", aggs["x"])

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	// get bins
	buckets, err := b.Bivariate.GetBins(coord, &res.Aggregations)
	if err != nil {
		return nil, err
	}

	// convert hit bins
	bins := make([][]map[string]interface{}, len(buckets))
	for i, bucket := range buckets {
		if bucket != nil {
			hits, err := b.TopHits.GetTopHits(&bucket.Aggregations)
			if err != nil {
				return nil, err
			}
			bins[i] = hits
		}
	}

	return json.Marshal(bins)
}
