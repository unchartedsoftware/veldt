package elastic

import (
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
	"github.com/unchartedsoftware/veldt/vendor/gopkg.in/olivere/elastic.v3"
)

// BivariateLPair represents an elasticsearch implementation of the bivariate tile.
type BivariatePair struct {
	Bivariate // elastic
	tile.BivariatePair
}

// GetQuery returns the tiling query.
func (b *BivariatePair) GetQuery(coord *binning.TileCoord) elastic.Query {

	fmt.Printf("<><><> elastic.bivariate_pair coord = %d %d %d \n", coord.Z, coord.X, coord.Y)

	// compute the tiling properties
	b.computeTilingProps(coord)

	fmt.Printf("<><><> elastic.bivariate_pair b.maxX: %d \n", b.Bivariate.maxX)

	// create the range queries
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewRangeQuery(b.BivariatePair.Bivariate.XField).
		Gte(b.Bivariate.minX).
		Lt(b.Bivariate.maxX))
	query.Must(elastic.NewRangeQuery(b.BivariatePair.Bivariate.YField).
		Gte(b.Bivariate.minY).
		Lt(b.Bivariate.maxY))
	return query
}

func (b *BivariatePair) computeTilingProps(coord *binning.TileCoord) {
	if b.tiling {
		return
	}
	// tiling params
	extents := &binning.Bounds{
		BottomLeft: &binning.Coord{
			X: b.BivariatePair.Left,
			Y: b.BivariatePair.Bottom,
		},
		TopRight: &binning.Coord{
			X: b.BivariatePair.Right,
			Y: b.BivariatePair.Top,
		},
	}
	fmt.Printf("<><><> elastic.bivariate_pair binning.Bounds = %f %f %f %f \n", b.BivariatePair.Left, b.BivariatePair.Right, b.BivariatePair.Bottom, b.BivariatePair.Top)
	b.BivariatePair.Bounds = binning.GetTileBounds(coord, extents)

	b.minX = int64(math.Min(b.BivariatePair.Bounds.BottomLeft.X, b.BivariatePair.Bounds.TopRight.X))
	b.maxX = int64(math.Max(b.BivariatePair.Bounds.BottomLeft.X, b.BivariatePair.Bounds.TopRight.X))
	b.minY = int64(math.Min(b.BivariatePair.Bounds.BottomLeft.Y, b.BivariatePair.Bounds.TopRight.Y))
	b.maxY = int64(math.Max(b.BivariatePair.Bounds.BottomLeft.Y, b.BivariatePair.Bounds.TopRight.Y))

	fmt.Printf("<><><> elastic.bivariate_pair Bounds = %d %d %d %d \n", b.minX, b.maxX, b.minY, b.maxY)
	// flag as computed
	b.tiling = true
}
