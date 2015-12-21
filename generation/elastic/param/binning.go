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
	Resolution int64
	BinSizeX   int64
	BinSizeY   int64
}

// NewBinning instantiates and returns a new binning parameter object.
func NewBinning(tileReq *tile.Request) (*Binning, error) {
	params := tileReq.Params
	tiling, err := NewTiling(tileReq)
	if err != nil {
		return nil, err
	}
	resolution := int64(json.GetNumberDefault(params, "resolution", binning.MaxTileResolution))
	return &Binning{
		Tiling:     tiling,
		BinSizeX:   int64(math.Abs(float64(tiling.Right-tiling.Left))) / resolution,
		BinSizeY:   int64(math.Abs(float64(tiling.Bottom-tiling.Top))) / resolution,
		Resolution: resolution,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Binning) GetHash() string {
	return fmt.Sprintf("%s:%d",
		p.Tiling.GetHash(),
		p.Resolution)
}

// GetXAgg returns an elastic aggregation.
func (p *Binning) GetXAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(p.Tiling.X).
		Interval(p.BinSizeX).
		MinDocCount(1)
}

// GetYAgg returns an elastic aggregation.
func (p *Binning) GetYAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(p.Tiling.Y).
		Interval(p.BinSizeY).
		MinDocCount(1)
}
