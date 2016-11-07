package param

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/param"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	param.Bivariate
}

// NewBivariate sets the params for the specific generator.
func NewBivariate(coord *binning.TileCoord, params map[string]interface{}) (*Bivariate, error) {
	// get x and y fields
	xField, ok := json.GetString(params, "xField")
	if !ok {
		return fmt.Errorf("`xField` parameter missing from tiling params")
	}
	yField, ok := json.GetString(params, "yField")
	if !ok {
		return fmt.Errorf("`yField` parameter missing from tiling params")
	}
	// get left, right, bottom, top extrema
	left, ok := json.GetNumber(params, "left")
	if !ok {
		return fmt.Errorf("`left` parameter missing from tiling params")
	}
	right, ok := json.GetNumber(params, "right")
	if !ok {
		return fmt.Errorf("`right` parameter missing from tiling params")
	}
	bottom, ok := json.GetNumber(params, "bottom")
	if !ok {
		return fmt.Errorf("`bottom` parameter missing from tiling params")
	}
	top, ok := json.GetNumber(params, "top")
	if !ok {
		return fmt.Errorf("`top` parameter missing from tiling params")
	}
	// get resolution
	resolution := json.GetNumberDefault(params, 256, "resolution")
	// get the tiles bounds
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: left,
			Y: top,
		},
		BottomRight: &binning.Coord{
			X: right,
			Y: bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	// get bin size
	xRange := math.Abs(bounds.BottomRight.X - bounds.TopLeft.X)
	yRange := math.Abs(bounds.BottomRight.Y - bounds.TopLeft.Y)
	binSizeX := xRange / resolution
	binSizeY := yRange / resolution
	// create bivariate
	b := &Bivariate{}
	b.XField = xField
	b.YField = yField
	b.Bounds = bounds
	b.MinX = math.Min(bounds.TopLeft.X, bounds.BottomRight.X)
	b.MaxX = math.Max(bounds.TopLeft.X, bounds.BottomRight.X)
	b.MinY = math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y)
	b.MaxY = math.Max(bounds.TopLeft.Y, bounds.BottomRight.Y)
	b.XRange = xRange
	b.YRange = yRange
	// add binning params
	b.Resolution = int(resolution)
	b.BinSizeX = binSizeX
	b.BinSizeY = binSizeY
	return b, nil
}

// GetXQuery returns the x query.
func (p *Bivariate) GetXQuery() elastic.Query {
	return elastic.NewRangeQuery(p.XField).
		Gte(p.MinX).
		Lt(p.MaxX)
}

// GetYQuery returns the y query.
func (p *Bivariate) GetYQuery() elastic.Query {
	return elastic.NewRangeQuery(p.YField).
		Gte(p.MinY).
		Lt(p.MaxY)
}

// ApplyXYAgg applies an elastic agg.
func (p *Bivariate) ApplyXYAgg(searchService *elastic.SearchService) elastic.Aggregation {
	intervalX := int64(math.Max(1, p.XRange/p.Resolution))
	intervalY := int64(math.Max(1, p.YRange/p.Resolution))
	x := elastic.NewHistogramAggregation().
			Field(p.XField).
			Offset(p.MinX).
			Interval(intervalX).
			MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(p.YField).
		Offset(p.MinY).
		Interval(intervalY).
		MinDocCount(1)
	x.SubAggregation("y", yAgg)
	return search.Aggregation("x", x)
}

// GetXYBins parses the resulting histograms into bins.
func (g *Bivariate) GetXYBins(res *elastic.SearchResult) ([]*elastic.AggregationBucketHistogramItem, error) {
	// parse aggregations
	xAgg, ok := res.Aggregations.Histogram("x")
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation `x` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, p.Resolution*p.Resolution)
	// fill bins
	for xBin, xBucket := range xAgg.Buckets {
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, fmt.Errorf("Histogram aggregation `y` was not found")
		}
		for yBin, yBucket := range yAgg.Buckets {
			index := xBin + p.Resolution*yBin
			bins[index] = yBucket
		}
	}
	return bins, nil
}
