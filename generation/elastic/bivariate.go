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
}

// GetQuery returns the tiling query.
func (b *Bivariate) GetQuery(coord *binning.TileCoord) elastic.Query {
	// get tile bounds
	bounds := b.TileBounds(coord)
	// create the range queries
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewRangeQuery(b.XField).
		Gte(int64(bounds.MinX())).
		Lt(int64(bounds.MaxX())))
	query.Must(elastic.NewRangeQuery(b.YField).
		Gte(int64(bounds.MinY())).
		Lt(int64(bounds.MaxY())))
	return query
}

// GetAggs returns the tiling aggregation.
func (b *Bivariate) GetAggs(coord *binning.TileCoord) map[string]elastic.Aggregation {
	bounds := b.TileBounds(coord)
	// compute binning itnernal
	intervalX := int64(math.Max(1, b.BinSizeX(coord)))
	intervalY := int64(math.Max(1, b.BinSizeY(coord)))
	// create the binning aggregations
	x := elastic.NewHistogramAggregation().
		Field(b.XField).
		Offset(int64(bounds.MinX())).
		Interval(intervalX).
		MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(b.YField).
		Offset(int64(bounds.MinY())).
		Interval(intervalY).
		MinDocCount(1)
	x.SubAggregation("y", y)
	return map[string]elastic.Aggregation{
		"x": x,
		"y": y,
	}
}

// GetAggsWithNested returns the tiling aggregation with a nested child agg.
func (b *Bivariate) GetAggsWithNested(coord *binning.TileCoord, id string, nested elastic.Aggregation) map[string]elastic.Aggregation {
	bounds := b.TileBounds(coord)
	// compute binning itnernal
	intervalX := int64(math.Max(1, b.BinSizeX(coord)))
	intervalY := int64(math.Max(1, b.BinSizeY(coord)))
	// create the binning aggregations
	x := elastic.NewHistogramAggregation().
		Field(b.XField).
		Offset(int64(bounds.MinX())).
		Interval(intervalX).
		MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(b.YField).
		Offset(int64(bounds.MinY())).
		Interval(intervalY).
		MinDocCount(1)
	x.SubAggregation("y", y)
	aggs := map[string]elastic.Aggregation{
		"x": x,
		"y": y,
	}
	if nested != nil {
		y.SubAggregation(id, nested)
		aggs[id] = nested
	}
	return aggs
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(coord *binning.TileCoord, aggs *elastic.Aggregations) ([]*elastic.AggregationBucketHistogramItem, error) {
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
		xBin := b.GetXBin(coord, float64(x))
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, fmt.Errorf("histogram aggregation `y` was not found")
		}
		for _, yBucket := range yAgg.Buckets {
			y := yBucket.Key
			yBin := b.GetYBin(coord, float64(y))
			index := xBin + b.Resolution*yBin
			bins[index] = yBucket
		}
	}
	return bins, nil
}
