package param

import (
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// Binning represents params for binning the data within the tile.
type Binning struct {
	Tiling     *Tiling
	Z          string
	Metric     string
	Resolution int64
	BinSizeX   float64
	BinSizeY   float64
	xInterval  int64
	yInterval  int64
}

// NewBinning instantiates and returns a new binning parameter object.
func NewBinning(tileReq *tile.Request) (*Binning, error) {
	params := tileReq.Params
	tiling, err := NewTiling(tileReq)
	if err != nil {
		return nil, err
	}
	bounds := tiling.Bounds
	xRange := math.Abs(bounds.BottomRight.X - bounds.TopLeft.X)
	yRange := math.Abs(bounds.BottomRight.Y - bounds.TopLeft.Y)
	resolution := json.GetNumberDefault(params, "resolution", binning.MaxTileResolution)
	binSizeX := xRange / resolution
	binSizeY := yRange / resolution
	return &Binning{
		Tiling:     tiling,
		Z:          json.GetStringDefault(params, "z", ""),
		Metric:     json.GetStringDefault(params, "metric", "sum"),
		Resolution: int64(resolution),
		BinSizeX:   binSizeX,
		BinSizeY:   binSizeY,
		xInterval:  int64(binSizeX),
		yInterval:  int64(binSizeY),
	}, nil
}

/*
func (p *Binning) clampBinToResolution(bin int64) int64 {
	if bin > p.Resolution-1 {
		return p.Resolution-1
	}
	if bin < 0 {
		return 0
	}
	return bin
}
*/

// GetHash returns a string hash of the parameter state.
func (p *Binning) GetHash() string {
	return fmt.Sprintf("%s:%s:%s:%d",
		p.Tiling.GetHash(),
		p.Z,
		p.Metric,
		p.Resolution)
}

// GetXAgg returns an elastic aggregation.
func (p *Binning) GetXAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(p.Tiling.X).
		Interval(p.xInterval).
		MinDocCount(1)
}

// GetYAgg returns an elastic aggregation.
func (p *Binning) GetYAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(p.Tiling.Y).
		Interval(p.yInterval).
		MinDocCount(1)
}

// GetZAgg returns an elastic aggregation.
func (p *Binning) GetZAgg() elastic.Aggregation {
	switch p.Metric {
	case "min":
		return elastic.NewMinAggregation().
			Field(p.Z)
	case "max":
		return elastic.NewMaxAggregation().
			Field(p.Z)
	case "avg":
		return elastic.NewAvgAggregation().
			Field(p.Z)
	default:
		return elastic.NewSumAggregation().
			Field(p.Z)
	}
}

// GetXBin given an x value, returns the corresponding bin.
func (p *Binning) GetXBin(x int64) int64 {
	bounds := p.Tiling.Bounds
	fx := float64(x)
	var bin int64
	if bounds.TopLeft.X > bounds.BottomRight.X {
		maxBin := float64(p.Resolution - 1)
		bin = int64(maxBin - ((fx - bounds.BottomRight.X) / p.BinSizeX))
	} else {
		bin = int64((fx - bounds.TopLeft.X) / p.BinSizeX)
	}
	return bin
}

// GetYBin given an y value, returns the corresponding bin.
func (p *Binning) GetYBin(y int64) int64 {
	bounds := p.Tiling.Bounds
	fy := float64(y)
	var bin int64
	if bounds.TopLeft.Y > bounds.BottomRight.Y {
		maxBin := float64(p.Resolution - 1)
		bin = int64(maxBin - ((fy - bounds.BottomRight.Y) / p.BinSizeY))
	} else {
		bin = int64((fy - bounds.TopLeft.Y) / p.BinSizeY)
	}
	return bin
}

// GetZAggValue extracts the value metric based on the type of operation
// specified.
func (p *Binning) GetZAggValue(aggName string, aggs *elastic.AggregationBucketHistogramItem) (*elastic.AggregationValueMetric, bool) {
	switch p.Metric {
	case "min":
		return aggs.Min(aggName)
	case "max":
		return aggs.Max(aggName)
	case "avg":
		return aggs.Avg(aggName)
	default:
		return aggs.Sum(aggName)
	}
}
