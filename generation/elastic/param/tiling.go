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
	defaultXField        = "pixel.x"
	defaultYField        = "pixel.y"
	defaultXType         = ""
	defaultXRelationship = ""
	defaultYType         = ""
	defaultYRelationship = ""
	defaultLeft          = 0.0
	defaultTop           = 0.0
	defaultRight         = binning.MaxPixels
	defaultBottom        = binning.MaxPixels
)

// Tiling represents params for tiling data.
type Tiling struct {
	X             string
	Y             string
	xType         string
	xRelationship string
	yType         string
	yRelationship string
	Bounds        *binning.Bounds
	minX          int64
	maxX          int64
	minY          int64
	maxY          int64
}

// NewTiling instantiates and returns a new tiling parameter object.
func NewTiling(tileReq *tile.Request) (*Tiling, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "binning")
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: json.GetNumberDefault(params, defaultLeft, "left"),
			Y: json.GetNumberDefault(params, defaultTop, "top"),
		},
		BottomRight: &binning.Coord{
			X: json.GetNumberDefault(params, defaultRight, "right"),
			Y: json.GetNumberDefault(params, defaultBottom, "bottom"),
		},
	}
	bounds := binning.GetTileBounds(tileReq.Coord, extents)
	return &Tiling{
		X:             json.GetStringDefault(params, defaultXField, "x"),
		Y:             json.GetStringDefault(params, defaultYField, "y"),
		xType:         json.GetStringDefault(params, defaultXType, "xType"),
		xRelationship: json.GetStringDefault(params, defaultXRelationship, "xRelationship"),
		yType:         json.GetStringDefault(params, defaultYType, "yType"),
		yRelationship: json.GetStringDefault(params, defaultYRelationship, "yRelationship"),
		Bounds:        bounds,
		minX:          int64(math.Min(bounds.TopLeft.X, bounds.BottomRight.X)),
		maxX:          int64(math.Max(bounds.TopLeft.X, bounds.BottomRight.X)),
		minY:          int64(math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y)),
		maxY:          int64(math.Max(bounds.TopLeft.Y, bounds.BottomRight.Y)),
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
func (p *Tiling) GetXQuery() elastic.Query {
	rangeQuery := elastic.NewRangeQuery(p.X).
		Gte(p.minX).
		Lt(p.maxX)
	if p.xType == "" || p.xRelationship == "" {
		return rangeQuery
	}
	if p.xRelationship == "child" {
		return elastic.NewHasChildQuery(
			p.xType,
			elastic.NewBoolQuery().Must(rangeQuery))
	}
	return elastic.NewHasParentQuery(
		p.xType,
		elastic.NewBoolQuery().Must(rangeQuery))
}

// GetYQuery returns an elastic query.
func (p *Tiling) GetYQuery() elastic.Query {
	rangeQuery := elastic.NewRangeQuery(p.Y).
		Gte(p.minY).
		Lt(p.maxY)
	if p.yType == "" || p.yRelationship == "" {
		return rangeQuery
	}
	if p.yRelationship == "child" {
		return elastic.NewHasChildQuery(
			p.yType,
			elastic.NewBoolQuery().Must(rangeQuery))
	}
	return elastic.NewHasParentQuery(
		p.yType,
		elastic.NewBoolQuery().Must(rangeQuery))
}
