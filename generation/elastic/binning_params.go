package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// BinningParams represents params for binning the data within the tile.
type BinningParams struct {
	X          string
	Y          string
	MinX       int64
	MaxX       int64
	MinY       int64
	MaxY       int64
	Resolution int64
	BinSizeX   int64
	BinSizeY   int64
}

// NewBinningParams parses the params map returns a pointer to the param struct.
func NewBinningParams(tileReq *tile.Request) *BinningParams {
	params := tileReq.Params
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: json.GetNumberDefault(params, "minX", 0.0),
			Y: json.GetNumberDefault(params, "maxX", binning.MaxPixels),
		},
		BottomRight: &binning.Coord{
			X: json.GetNumberDefault(params, "maxY", 0.0),
			Y: json.GetNumberDefault(params, "minY", binning.MaxPixels),
		},
	}
	bounds := binning.GetTileBounds(tileReq.TileCoord, extents)
	resolution := int64(json.GetNumberDefault(params, "resolution", binning.MaxTileResolution))
	return &BinningParams{
		X: json.GetStringDefault(params, "x", "pixel.x"),
		Y: json.GetStringDefault(params, "y", "pixel.y"),
		BinSizeX: int64(bounds.BottomRight.X-bounds.TopLeft.X) / resolution,
		BinSizeY: int64(bounds.BottomRight.Y-bounds.TopLeft.Y) / resolution,
		MinX: int64(bounds.TopLeft.X),
		MaxX: int64(bounds.BottomRight.X - 1),
		MinY: int64(bounds.TopLeft.Y),
		MaxY: int64(bounds.BottomRight.Y - 1),
		Resolution: resolution,
	}
}

// GetHash returns a string hash of the parameter state.
func (p *BinningParams) GetHash() string {
	return fmt.Sprintf("%s:%s:%d:%d:%d:%d:%d",
		p.X,
		p.Y,
		p.MinX,
		p.MaxX,
		p.MinY,
		p.MaxY,
		p.Resolution)
}

// GetXQuery returns an elastic query.
func (p *BinningParams) GetXQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.X).
		Gte(p.MinX).
		Lte(p.MaxX)
}

// GetYQuery returns an elastic query.
func (p *BinningParams) GetYQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.Y).
		Gte(p.MinY).
		Lte(p.MaxY)
}

// GetXAgg returns an elastic aggregation.
func (p *BinningParams) GetXAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(p.X).
		Interval(p.BinSizeX).
		MinDocCount(1)
}

// GetYAgg returns an elastic aggregation.
func (p *BinningParams) GetYAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(p.Y).
		Interval(p.BinSizeY).
		MinDocCount(1)
}
