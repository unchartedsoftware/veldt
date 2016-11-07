package param

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/param"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// Univariate represents a univariate tile generator.
type Univariate struct {
	param.Univariate
}

// NewUnivariate sets the params for the specific generator.
func NewUnivariate(coord *binning.TileCoord, params map[string]interface{}) (*Univariate, error) {
	// get x and y fields
	field, ok := json.GetString(params, "field")
	if !ok {
		return fmt.Errorf("`field` parameter missing from tiling params")
	}
	// get left, right, bottom, top extrema
	min, ok := json.GetNumber(params, "min")
	if !ok {
		return fmt.Errorf("`min` parameter missing from tiling params")
	}
	max, ok := json.GetNumber(params, "max")
	if !ok {
		return fmt.Errorf("`max` parameter missing from tiling params")
	}
	extrema := binning.GetTileExtrema(coord.X, coord.Z, &binning.Extrema{
		Min: min,
		Max: max,
	})
	// get resolution
	resolution := json.GetNumberDefault(params, 256, "resolution")
	// get bin size
	rang := math.Abs(extrema.Min - extrema.Max)
	binSize := rang / resolution
	// create univariate
	u := &Univariate{}
	u.Field = field
	u.Min = extrema.Min
	u.Max = extrema.Max
	u.Range = rang
	// add binning params
	u.Resolution = int(resolution)
	u.BinSize = binSize
	return u, nil
}

// GetQuery returns the x query.
func (p *Bivariate) GetQuery() elastic.Query {
	return elastic.NewRangeQuery(p.Field).
		Gte(p.Min).
		Lt(p.Max)
}

// ApplyAgg applies an elastic agg.
func (p *Bivariate) ApplyAgg(searchService *elastic.SearchService) elastic.Aggregation {
	interval := int64(math.Max(1, p.Range/p.Resolution))
	agg := elastic.NewHistogramAggregation().
			Field(p.Field).
			Offset(p.Min).
			Interval(interval).
			MinDocCount(1)
	return search.Aggregation("bins", agg)
}

// GetBins parses the resulting histograms into bins.
func (g *Bivariate) GetBins(res *elastic.SearchResult) ([]*elastic.AggregationBucketHistogramItem, error) {
	// parse aggregations
	agg, ok := res.Aggregations.Histogram("bins")
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation `bins` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, p.Resolution)
	// fill bins buffer
	for i, bucket := range agg.Buckets {
		bins[i] = bucket
	}
	return bins, nil
}
