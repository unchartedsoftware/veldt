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
	tiling bool
	minX   int64
	maxX   int64
	minY   int64
	maxY   int64
	// binning
	binning   bool
	intervalX int64
	intervalY int64
}

func (b *Bivariate) computeTilingProps(coord *binning.TileCoord) {
	if b.tiling {
		return
	}
	// tiling params
	extents := &binning.Bounds{
		BottomLeft: &binning.Coord{
			X: b.Left,
			Y: b.Bottom,
		},
		TopRight: &binning.Coord{
			X: b.Right,
			Y: b.Top,
		},
	}
	b.Bounds = binning.GetTileBounds(coord, extents)
	b.minX = int64(math.Min(b.Bounds.BottomLeft.X, b.Bounds.TopRight.X))
	b.maxX = int64(math.Max(b.Bounds.BottomLeft.X, b.Bounds.TopRight.X))
	b.minY = int64(math.Min(b.Bounds.BottomLeft.Y, b.Bounds.TopRight.Y))
	b.maxY = int64(math.Max(b.Bounds.BottomLeft.Y, b.Bounds.TopRight.Y))
	// flag as computed
	b.tiling = true
}

func (b *Bivariate) computeBinningProps(coord *binning.TileCoord) {
	if b.binning {
		return
	}
	// ensure we have tiling props
	b.computeTilingProps(coord)
	// binning params
	xRange := math.Abs(b.Bounds.TopRight.X - b.Bounds.BottomLeft.X)
	yRange := math.Abs(b.Bounds.TopRight.Y - b.Bounds.BottomLeft.Y)
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
		Gte(b.minX).
		Lt(b.maxX))
	query.Must(elastic.NewRangeQuery(b.YField).
		Gte(b.minY).
		Lt(b.maxY))
	return query
}

// GetAggs returns the tiling aggregation.
func (b *Bivariate) GetAggs(coord *binning.TileCoord) map[string]elastic.Aggregation {
	// compute the binning properties
	b.computeBinningProps(coord)
	// create the binning aggregations
	x := elastic.NewHistogramAggregation().
		Field(b.XField).
		Offset(b.minX).
		Interval(b.intervalX).
		MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(b.YField).
		Offset(b.minY).
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
