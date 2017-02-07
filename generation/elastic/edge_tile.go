package elastic

import (
	"fmt"
	"math"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
)

// EdgeTile represents an elasticsearch implementation of the Edge tile.
type EdgeTile struct {
	Tile
	TopHits
	tile.Edge
	// Extremes of the single tile
	minX             int64
	maxX             int64
	minY             int64
	maxY             int64
	isTilingComputed bool
}

// NewEdgeTile instantiates and returns a new tile struct.
func NewEdgeTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		e := &EdgeTile{}
		e.Host = host
		e.Port = port
		return e, nil
	}
}

//TODO: move to edge.go. this isn't elastic-specific afaik.
func (e *EdgeTile) computeTilingProps(coord *binning.TileCoord) {
	if e.isTilingComputed {
		return
	}
	// tiling params
	e.TileBounds = binning.GetTileBounds(coord, e.WorldBounds)
	e.minX = int64(math.Min(e.TileBounds.BottomLeft().X, e.TileBounds.TopRight().X))
	e.maxX = int64(math.Max(e.TileBounds.BottomLeft().X, e.TileBounds.TopRight().X))
	e.minY = int64(math.Min(e.TileBounds.BottomLeft().Y, e.TileBounds.TopRight().Y))
	e.maxY = int64(math.Max(e.TileBounds.BottomLeft().Y, e.TileBounds.TopRight().Y))
	// flag as computed
	e.isTilingComputed = true
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (e *EdgeTile) Parse(params map[string]interface{}) error {

	err := e.Edge.Parse(params)
	if err != nil {
		return err
	}

	err = e.TopHits.Parse(params)
	if err != nil {
		return err
	}

	// parse includes
	e.TopHits.IncludeFields = e.ParseIncludes(e.TopHits.IncludeFields)

	return nil
}

// GetQuery returns the tiling query.
func (e *EdgeTile) GetQuery(coord *binning.TileCoord) elastic.Query {
	e.computeTilingProps(coord)

	// create the range queries
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewRangeQuery(e.SrcXField).
		Gte(e.minX).
		Lt(e.maxX))
	query.Must(elastic.NewRangeQuery(e.SrcYField).
		Gte(e.minY).
		Lt(e.maxY))
	return query
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (e *EdgeTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// get client
	client, err := NewClient(e.Host, e.Port)
	if err != nil {
		return nil, err
	}
	// create search service
	search := client.Search().
		Index(uri).
		Size(0)

	// create root query
	q, err := e.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q = q.Must(e.GetQuery(coord))

	// set the query
	search.Query(q)

	// get aggs
	aggs := e.TopHits.GetAggs()
	// set the aggregation
	search.Aggregation("top-hits", aggs["top-hits"])

	// send query
	res, err := search.Pretty(true).Do()
	if err != nil {
		_ = fmt.Errorf("EdgeTile: query error %s\n", err)
		return nil, err
	}

	// get top hits
	hits, err := e.TopHits.GetTopHits(&res.Aggregations)
	if err != nil {
		return nil, err
	}

	// convert to point array
	points := make([]float32, len(hits)*4)
	for i, hit := range hits {

		zoom := uint32(0)

		// get hit x/y in tile coords
		x, y, ok := e.GetSrcXY(hit, zoom)
		if !ok {
			_ = fmt.Errorf("Couldn't GetSrcXY\n")
			continue
		}
		x2, y2, ok := e.GetDstXY(hit, zoom)
		if !ok {
			_ = fmt.Errorf("Couldn't GetDstXY\n")
			continue
		}
		// add to point array
		points[i*4] = x
		points[i*4+1] = y
		points[i*4+2] = x2
		points[i*4+3] = y2
	}
	// encode and return results
	return e.Edge.Encode(hits, points)
}
