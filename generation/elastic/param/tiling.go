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
	X    string
	Y    string
	MinX int64
	MaxX int64
	MinY int64
	MaxY int64
}

// NewTiling instantiates and returns a new tiling parameter object.
func NewTiling(tileReq *tile.Request) (*Tiling, error) {
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
	return &Tiling{
		X:    json.GetStringDefault(params, "x", "pixel.x"),
		Y:    json.GetStringDefault(params, "y", "pixel.y"),
		MinX: int64(bounds.TopLeft.X),
		MaxX: int64(bounds.BottomRight.X),
		MinY: int64(bounds.TopLeft.Y),
		MaxY: int64(bounds.BottomRight.Y),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Tiling) GetHash() string {
	return fmt.Sprintf("%s:%s:%d:%d:%d:%d",
		p.X,
		p.Y,
		p.MinX,
		p.MaxX,
		p.MinY,
		p.MaxY)
}

// GetXQuery returns an elastic query.
func (p *Tiling) GetXQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.X).
		Gte(p.MinX).
		Lt(p.MaxX)
}

// GetYQuery returns an elastic query.
func (p *Tiling) GetYQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.Y).
		Gte(p.MinY).
		Lt(p.MaxY)
}
