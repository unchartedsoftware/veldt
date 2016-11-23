package elastic

import (
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	tile.Bivariate
	// tiling
	tiling bool
	bounds *binning.Bounds
	minX   int64
	maxX   int64
	minY   int64
	maxY   int64
	// binning
	binning   bool
	binSizeX  float64
	binSizeY  float64
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
	b.bounds = binning.GetTileBounds(coord, extents)
	b.minX = int64(math.Min(b.bounds.BottomLeft.X, b.bounds.TopRight.X))
	b.maxX = int64(math.Max(b.bounds.BottomLeft.X, b.bounds.TopRight.X))
	b.minY = int64(math.Min(b.bounds.BottomLeft.Y, b.bounds.TopRight.Y))
	b.maxY = int64(math.Max(b.bounds.BottomLeft.Y, b.bounds.TopRight.Y))
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
	xRange := math.Abs(b.bounds.TopRight.X - b.bounds.BottomLeft.X)
	yRange := math.Abs(b.bounds.TopRight.Y - b.bounds.BottomLeft.Y)
	b.intervalX = int64(math.Max(1, xRange/float64(b.Resolution)))
	b.intervalY = int64(math.Max(1, yRange/float64(b.Resolution)))
	b.binSizeX = xRange / float64(b.Resolution)
	b.binSizeY = yRange / float64(b.Resolution)
	// flag as computed
	b.binning = true
}

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

// GetXBin given an x value, returns the corresponding bin.
func (b *Bivariate) getXBin(x int64) int {
	bounds := b.bounds
	fx := float64(x)
	var bin int64
	if bounds.BottomLeft.X > bounds.TopRight.X {
		bin = int64(float64(b.Resolution-1) - ((fx - bounds.TopRight.X) / b.binSizeX))
	} else {
		bin = int64((fx - bounds.BottomLeft.X) / b.binSizeX)
	}
	return b.clampBin(bin)
}

// GetYBin given an y value, returns the corresponding bin.
func (b *Bivariate) getYBin(y int64) int {
	bounds := b.bounds
	fy := float64(y)
	var bin int64
	if bounds.BottomLeft.Y > bounds.TopRight.Y {
		bin = int64(float64(b.Resolution-1) - ((fy - bounds.TopRight.Y) / b.binSizeY))
	} else {
		bin = int64((fy - bounds.BottomLeft.Y) / b.binSizeY)
	}
	return b.clampBin(bin)
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(res *elastic.SearchResult) ([]*elastic.AggregationBucketHistogramItem, error) {
	if !b.binning {
		return nil, fmt.Errorf("binning properties have not been computed, ensure `GetAggs` is called")
	}
	// parse aggregations
	xAgg, ok := res.Aggregations.Histogram("x")
	if !ok {
		return nil, fmt.Errorf("histogram aggregation `x` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, b.Resolution*b.Resolution)
	// fill bins
	for _, xBucket := range xAgg.Buckets {
		x := xBucket.Key
		xBin := b.getXBin(x)
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, fmt.Errorf("histogram aggregation `y` was not found")
		}
		for _, yBucket := range yAgg.Buckets {
			y := yBucket.Key
			yBin := b.getYBin(y)
			index := xBin + b.Resolution*yBin
			bins[index] = yBucket
		}
	}
	return bins, nil
}

func (b *Bivariate) clampBin(bin int64) int {
	if bin > int64(b.Resolution)-1 {
		return b.Resolution - 1
	}
	if bin < 0 {
		return 0
	}
	return int(bin)
}
