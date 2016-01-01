package param

import (
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// Tiling represents params for tiling data.
type Tiling struct {
	X      string
	Y      string
	Bounds *binning.Bounds
	minX   int64
	maxX   int64
	minY   int64
	maxY   int64
}

// NewTiling instantiates and returns a new tiling parameter object.
func NewTiling(tileReq *tile.Request) (*Tiling, error) {
	params := tileReq.Params
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: json.GetNumberDefault(params, "left", 0.0),
			Y: json.GetNumberDefault(params, "top", 0.0),
		},
		BottomRight: &binning.Coord{
			X: json.GetNumberDefault(params, "right", binning.MaxPixels),
			Y: json.GetNumberDefault(params, "bottom", binning.MaxPixels),
		},
	}
	bounds := binning.GetTileBounds(tileReq.Coord, extents)
	return &Tiling{
		X:      json.GetStringDefault(params, "x", "pixel.x"),
		Y:      json.GetStringDefault(params, "y", "pixel.y"),
		Bounds: bounds,
		minX:   int64(math.Min(bounds.TopLeft.X, bounds.BottomRight.X)),
		maxX:   int64(math.Max(bounds.TopLeft.X, bounds.BottomRight.X)),
		minY:   int64(math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y)),
		maxY:   int64(math.Max(bounds.TopLeft.Y, bounds.BottomRight.Y)),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Tiling) GetHash() string {
	return fmt.Sprintf("%s:%s:%d:%d:%d:%d",
		p.X,
		p.Y,
		p.minX,
		p.maxX,
		p.minY,
		p.maxY)
}

// GetXQuery returns an elastic query.
func (p *Tiling) GetXQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.X).
		Gte(p.minX).
		Lt(p.maxX)
}

// GetYQuery returns an elastic query.
func (p *Tiling) GetYQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.Y).
		Gte(p.minY).
		Lt(p.maxY)
}
