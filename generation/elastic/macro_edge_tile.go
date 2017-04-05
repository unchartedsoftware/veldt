package elastic

import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
)

// MacroEdgeTile represents an elasticsearch implementation of the Edge tile.
type MacroEdgeTile struct {
	Tile
	TopHits
	Edge
	tile.MacroEdge
}

// NewMacroEdgeTile instantiates and returns a new tile struct.
func NewMacroEdgeTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		e := &MacroEdgeTile{}
		e.Host = host
		e.Port = port
		return e, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (e *MacroEdgeTile) Parse(params map[string]interface{}) error {
	err := e.Edge.Parse(params)
	if err != nil {
		return err
	}
	err = e.TopHits.Parse(params)
	if err != nil {
		return err
	}
	// parse includes
	e.TopHits.IncludeFields = e.MacroEdge.ParseIncludes(
		e.TopHits.IncludeFields,
		e.Edge.SrcXField,
		e.Edge.SrcYField,
		e.Edge.DstXField,
		e.Edge.DstYField)
	return e.MacroEdge.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (e *MacroEdgeTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
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
	q = q.Must(e.Edge.GetQuery(coord))

	// set the query
	search.Query(q)

	// get aggs
	aggs := e.TopHits.GetAggs()
	// set the aggregation
	search.Aggregation("top-hits", aggs["top-hits"])

	// send query
	res, err := search.Pretty(true).Do()
	if err != nil {
		return nil, err
	}

	// get top hits
	hits, err := e.TopHits.GetTopHits(&res.Aggregations)
	if err != nil {
		return nil, err
	}

	// convert to point array
	points := make([]float32, len(hits)*4)
	// get hit x/y in tile coords
	for i, hit := range hits {
		srcX, srcY, ok := e.Edge.GetSrcXY(coord, hit)
		if !ok {
			continue
		}
		dstX, dstY, ok := e.Edge.GetDstXY(coord, hit)
		if !ok {
			continue
		}
		// add to point array
		points[i*4] = float32(srcX)
		points[i*4+1] = float32(srcY)
		points[i*4+2] = float32(dstX)
		points[i*4+3] = float32(dstY)
	}

	// encode and return results
	return e.MacroEdge.Encode(points)
}
