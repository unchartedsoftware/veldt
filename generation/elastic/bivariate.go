package elastic

import (
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
)

// Bivariate represents an elasticsearch implementation of the bivariate tile.
type Bivariate struct {
	tile.Bivariate
	// tiling
	isTilingComputed bool

	// binning
	binning   bool
	intervalX int64
	intervalY int64
}

func (b *Bivariate) computeTilingProps(coord *binning.TileCoord) {
	if b.isTilingComputed {
		return
	}
	// tiling params
	b.TileBounds = binning.GetTileBounds(coord, b.WorldBounds)

	// flag as computed
	b.isTilingComputed = true
}

func (b *Bivariate) computeBinningProps(coord *binning.TileCoord) {
	if b.binning {
		return
	}
	// ensure we have tiling props
	b.computeTilingProps(coord)
	// binning params
	xRange := math.Abs(b.TileBounds.TopRight().X - b.TileBounds.BottomLeft().X)
	yRange := math.Abs(b.TileBounds.TopRight().Y - b.TileBounds.BottomLeft().Y)
	b.intervalX = int64(math.Max(1, xRange/float64(b.Resolution)))
	b.intervalY = int64(math.Max(1, yRange/float64(b.Resolution)))
	b.BinSizeX = xRange / float64(b.Resolution)
	b.BinSizeY = yRange / float64(b.Resolution)
	// flag as computed
	b.binning = true
}

// GetQuery returns the tiling query.
func (b *Bivariate) GetQuery(coord *binning.TileCoord) elastic.Query {
	// compute the tiling properties
	b.computeTilingProps(coord)
	// create the range queries
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewRangeQuery(b.XField).
		Gte(int64(b.TileBounds.MinX())).
		Lt(int64(b.TileBounds.MaxX())))
	query.Must(elastic.NewRangeQuery(b.YField).
		Gte(int64(b.TileBounds.MinY())).
		Lt(int64(b.TileBounds.MaxY())))
	return query
}

// GetAggs returns the tiling aggregation.
func (b *Bivariate) GetAggs(coord *binning.TileCoord) map[string]elastic.Aggregation {
	// compute the binning properties
	b.computeBinningProps(coord)
	// create the binning aggregations
	x := elastic.NewHistogramAggregation().
		Field(b.XField).
		Offset(int64(b.TileBounds.MinX())).
		Interval(b.intervalX).
		MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(b.YField).
		Offset(int64(b.TileBounds.MinY())).
		Interval(b.intervalY).
		MinDocCount(1)
	x.SubAggregation("y", y)
	return map[string]elastic.Aggregation{
		"x": x,
		"y": y,
	}
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(aggs *elastic.Aggregations) ([]*elastic.AggregationBucketHistogramItem, error) {
	if !b.binning {
		return nil, fmt.Errorf("binning properties have not been computed, ensure `GetAggs` is called")
	}
	// parse aggregations
	xAgg, ok := aggs.Histogram("x")
	if !ok {
		return nil, fmt.Errorf("histogram aggregation `x` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, b.Resolution*b.Resolution)
	// fill bins
	for _, xBucket := range xAgg.Buckets {
		x := xBucket.Key
		xBin := b.GetXBin(x)
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, fmt.Errorf("histogram aggregation `y` was not found")
		}
		for _, yBucket := range yAgg.Buckets {
			y := yBucket.Key
			yBin := b.GetYBin(y)
			index := xBin + b.Resolution*yBin
			bins[index] = yBucket
		}
	}
	return bins, nil
}
