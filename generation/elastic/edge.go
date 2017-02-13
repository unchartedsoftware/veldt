package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
)

// Edge represents an elasticsearch implementation of the bivariate tile.
type Edge struct {
	tile.Edge
}

// GetQuery returns the tiling query.
func (e *Edge) GetQuery(coord *binning.TileCoord) elastic.Query {
	// get tile bounds
	bounds := e.TileBounds(coord)
	// create the range queries
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewRangeQuery(e.Edge.SrcXField).
		Gte(int64(bounds.MinX())).
		Lt(int64(bounds.MaxX())))
	query.Must(elastic.NewRangeQuery(e.Edge.SrcYField).
		Gte(int64(bounds.MinY())).
		Lt(int64(bounds.MaxY())))
	return query
}
