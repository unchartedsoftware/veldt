package param

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// Tiling represents params for tiling data.
type Tiling struct {
	X      string
	Y      string
	Left   int64
	Right  int64
	Top    int64
	Bottom int64
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
	bounds := binning.GetTileBounds(tileReq.TileCoord, extents)
	return &Tiling{
		X:      json.GetStringDefault(params, "x", "pixel.x"),
		Y:      json.GetStringDefault(params, "y", "pixel.y"),
		Left:   int64(bounds.TopLeft.X),
		Right:  int64(bounds.BottomRight.X),
		Top:    int64(bounds.TopLeft.Y),
		Bottom: int64(bounds.BottomRight.Y),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Tiling) GetHash() string {
	return fmt.Sprintf("%s:%s:%d:%d:%d:%d",
		p.X,
		p.Y,
		p.Left,
		p.Right,
		p.Top,
		p.Bottom)
}

// GetXQuery returns an elastic query.
func (p *Tiling) GetXQuery() *elastic.RangeQuery {
	if p.Right > p.Left {
		return elastic.NewRangeQuery(p.X).
			Gte(p.Left).
			Lt(p.Right)
	}
	return elastic.NewRangeQuery(p.X).
		Gte(p.Right).
		Lt(p.Left)
}

// GetYQuery returns an elastic query.
func (p *Tiling) GetYQuery() *elastic.RangeQuery {
	if p.Top > p.Bottom {
		return elastic.NewRangeQuery(p.Y).
			Gte(p.Bottom).
			Lt(p.Top)
	}
	return elastic.NewRangeQuery(p.Y).
		Gte(p.Top).
		Lt(p.Bottom)
}
