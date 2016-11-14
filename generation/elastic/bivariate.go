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
}

func (b *Bivariate) GetQuery(coord *binning.TileCoord) elastic.Query {
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: b.Left,
			Y: b.Top,
		},
		BottomRight: &binning.Coord{
			X: b.Right,
			Y: b.Bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	minX := math.Min(bounds.TopLeft.X, bounds.BottomRight.X)
	maxX := math.Max(bounds.TopLeft.X, bounds.BottomRight.X)
	minY := math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y)
	maxY := math.Max(bounds.TopLeft.Y, bounds.BottomRight.Y)

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewRangeQuery(b.XField).
		Gte(minX).
		Lt(maxX))
	query.Must(elastic.NewRangeQuery(b.YField).
		Gte(minY).
		Lt(maxY))
	return query
}

func (b *Bivariate) GetAggs(coord *binning.TileCoord) map[string]elastic.Aggregation {
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: b.Left,
			Y: b.Top,
		},
		BottomRight: &binning.Coord{
			X: b.Right,
			Y: b.Bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	minX := int64(math.Min(bounds.TopLeft.X, bounds.BottomRight.X))
	minY := int64(math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y))
	xRange := math.Abs(bounds.BottomRight.X - bounds.TopLeft.X)
	yRange := math.Abs(bounds.BottomRight.Y - bounds.TopLeft.Y)
	intervalX := int64(math.Max(1, xRange/float64(b.Resolution)))
	intervalY := int64(math.Max(1, yRange/float64(b.Resolution)))

	x := elastic.NewHistogramAggregation().
		Field(b.XField).
		Offset(minX).
		Interval(intervalX).
		MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(b.YField).
		Offset(minY).
		Interval(intervalY).
		MinDocCount(1)

	x.SubAggregation("y", y)
	return map[string]elastic.Aggregation{
		"x": x,
		"y": y,
	}
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(res *elastic.SearchResult) ([]*elastic.AggregationBucketHistogramItem, error) {
	// parse aggregations
	xAgg, ok := res.Aggregations.Histogram("x")
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation `x` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, b.Resolution*b.Resolution)
	// fill bins
	for xBin, xBucket := range xAgg.Buckets {
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, fmt.Errorf("Histogram aggregation `y` was not found")
		}
		for yBin, yBucket := range yAgg.Buckets {
			index := xBin + b.Resolution*yBin
			bins[index] = yBucket
		}
	}
	return bins, nil
}
