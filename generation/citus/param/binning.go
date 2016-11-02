package param

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/citus/query"
	"github.com/unchartedsoftware/prism/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultResolution  = binning.MaxTileResolution
	defaultZField      = ""
	defaultMetric      = "sum"
	intervalResolution = 256
)

// Binning represents params for binning the data within the tile.
type Binning struct {
	Tiling     *Tiling
	Z          string
	Metric     string
	Resolution int64
	BinSizeX   float64
	BinSizeY   float64
	intervalX  int64
	intervalY  int64
}

// NewBinning instantiates and returns a new binning parameter object.
func NewBinning(tileReq *tile.Request) (*Binning, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "binning")
	tiling, err := NewTiling(tileReq)
	if err != nil {
		return nil, err
	}
	bounds := tiling.Bounds
	xRange := math.Abs(bounds.BottomRight.X - bounds.TopLeft.X)
	yRange := math.Abs(bounds.BottomRight.Y - bounds.TopLeft.Y)
	resolution := json.GetNumberDefault(params, defaultResolution, "resolution")
	binSizeX := xRange / resolution
	binSizeY := yRange / resolution
	return &Binning{
		Tiling:     tiling,
		Z:          json.GetStringDefault(params, defaultZField, "z"),
		Metric:     json.GetStringDefault(params, defaultMetric, "metric"),
		Resolution: int64(resolution),
		intervalX:  int64(math.Max(1, xRange/intervalResolution)),
		intervalY:  int64(math.Max(1, yRange/intervalResolution)),
		BinSizeX:   binSizeX,
		BinSizeY:   binSizeY,
	}, nil
}

func (p *Binning) clampBin(bin int64) int64 {
	if bin > p.Resolution-1 {
		return p.Resolution - 1
	}
	if bin < 0 {
		return 0
	}
	return bin
}

// GetHash returns a string hash of the parameter state.
func (p *Binning) GetHash() string {
	return fmt.Sprintf("%s:%s:%s:%d",
		p.Tiling.GetHash(),
		p.Z,
		p.Metric,
		p.Resolution)
}

// AddXAgg adds a groupby clause to the query.
func (p *Binning) AddXAgg(query *query.Query) *query.Query {
	minXArg := query.AddParameter(p.Tiling.minX)
	intervalXArg := query.AddParameter(p.intervalX)
	queryString := fmt.Sprintf("((%s - %s) / %s * %s)", p.Tiling.X, minXArg, intervalXArg, intervalXArg)
	query.AddGroupByClause(queryString)
	query.AddField(fmt.Sprintf("%s + %s as x", minXArg, queryString))
	//TODO: Handle the MinDocCount.

	return query
}

// AddYAgg adds a groupby clause to the query.
func (p *Binning) AddYAgg(query *query.Query) *query.Query {
	minYArg := query.AddParameter(p.Tiling.minY)
	intervalYArg := query.AddParameter(p.intervalY)
	queryString := fmt.Sprintf("((%s - %s) / %s * %s)", p.Tiling.Y, minYArg, intervalYArg, intervalYArg)
	query.AddGroupByClause(queryString)
	query.AddField(fmt.Sprintf("%s + %s as y", minYArg, queryString))
	//TODO: Handle the MinDocCount.

	return query
}

// AddZAgg adds an aggregation to the query.
func (p *Binning) AddZAgg(query *query.Query) *query.Query {
	fieldString := ""
	switch p.Metric {
	case "min":
		fieldString = fmt.Sprintf("MIN(%s)", p.Z)
	case "max":
		fieldString = fmt.Sprintf("MAX(%s)", p.Z)
	case "avg":
		fieldString = fmt.Sprintf("AVG(%s)", p.Z)
	default:
		fieldString = fmt.Sprintf("SUM(%s)", p.Z)
	}
	query.AddField(fieldString)

	return query
}

// GetXBin given an x value, returns the corresponding bin.
func (p *Binning) GetXBin(x int64) int64 {
	bounds := p.Tiling.Bounds
	fx := float64(x)
	var bin int64
	if bounds.TopLeft.X > bounds.BottomRight.X {
		bin = int64(float64(p.Resolution) - ((fx - bounds.BottomRight.X) / p.BinSizeX))
	} else {
		bin = int64((fx - bounds.TopLeft.X) / p.BinSizeX)
	}
	return p.clampBin(bin)
}

// GetYBin given an y value, returns the corresponding bin.
func (p *Binning) GetYBin(y int64) int64 {
	bounds := p.Tiling.Bounds
	fy := float64(y)
	var bin int64
	if bounds.TopLeft.Y > bounds.BottomRight.Y {
		bin = int64(float64(p.Resolution) - ((fy - bounds.BottomRight.Y) / p.BinSizeY))
	} else {
		bin = int64((fy - bounds.TopLeft.Y) / p.BinSizeY)
	}
	return p.clampBin(bin)
}

//TODO: Not sure how this fits in yet.
// GetZAggValue extracts the value metric based on the type of operation
// specified.
//func (p *Binning) GetZAggValue(aggName string, aggs *elastic.AggregationBucketHistogramItem) (*elastic.AggregationValueMetric, bool) {
//	switch p.Metric {
//	case "min":
//		return aggs.Min(aggName)
//	case "max":
//		return aggs.Max(aggName)
//	case "avg":
//		return aggs.Avg(aggName)
//	default:
//		return aggs.Sum(aggName)
//	}
//}
