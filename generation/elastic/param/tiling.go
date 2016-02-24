package param

import (
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultXField = "pixel.x"
	defaultYField = "pixel.y"
	defaultLeft   = 0.0
	defaultTop    = 0.0
	defaultRight  = binning.MaxPixels
	defaultBottom = binning.MaxPixels
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
	params := json.GetChildOrEmpty(tileReq.Params, "binning")
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: json.GetNumberDefault(params, "left", defaultLeft),
			Y: json.GetNumberDefault(params, "top", defaultTop),
		},
		BottomRight: &binning.Coord{
			X: json.GetNumberDefault(params, "right", defaultRight),
			Y: json.GetNumberDefault(params, "bottom", defaultBottom),
		},
	}
	bounds := binning.GetTileBounds(tileReq.Coord, extents)
	return &Tiling{
		X:      json.GetStringDefault(params, "x", defaultXField),
		Y:      json.GetStringDefault(params, "y", defaultYField),
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
