package param

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/citus/query"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultXField        = "pixel_x"
	defaultYField        = "pixel_y"
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

// AddXQuery returns an elastic query.
func (p *Tiling) AddXQuery(query *query.Query) *query.Query {
	minXArg := query.AddParameter(p.minX)
	maxXArg := query.AddParameter(p.maxX)
	rangeQuery := fmt.Sprintf("%s >= %s and %s < %s", p.X, minXArg, p.X, maxXArg)
	query.AddWhereClause(rangeQuery)

	//TODO: Find a way to flag documents without values...perhaps using null
	query.AddWhereClause(fmt.Sprintf("%s != 0", p.X))

	return query
}

// AddYQuery returns an elastic query.
func (p *Tiling) AddYQuery(query *query.Query) *query.Query {
	minYArg := query.AddParameter(p.minY)
	maxYArg := query.AddParameter(p.maxY)
	rangeQuery := fmt.Sprintf("%s >= %s and %s < %s", p.Y, minYArg, p.Y, maxYArg)
	query.AddWhereClause(rangeQuery)

	//TODO: Find a way to flag documents without values...perhaps using null
	query.AddWhereClause(fmt.Sprintf("%s != 0", p.Y))

	return query
}
