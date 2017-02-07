package geometry

import (
	"fmt"

	jsonUtil "github.com/unchartedsoftware/veldt/util/json"
)

// Bounds represents an Immutable bounding rectangle with convenience getters.
type Bounds struct {
	left   float64
	right  float64
	bottom float64
	top    float64
}

// Rectangle represents a pair of Coords
type Rectangle struct {
	BottomLeft *Coord
	TopRight   *Coord
}

// NewBounds returns a new Bounds object
func NewBounds(left, right, bottom, top float64) *Bounds {
	return &Bounds{
		left:   left,
		right:  right,
		bottom: bottom,
		top:    top,
	}
}

// NewBoundsFromRectangle creates a Bounds from a Rectangle's corners
func NewBoundsFromRectangle(corners *Rectangle) *Bounds {
	return &Bounds{
		left:   corners.BottomLeft.X,
		right:  corners.TopRight.X,
		bottom: corners.BottomLeft.Y,
		top:    corners.TopRight.Y,
	}
}

// NewBoundsByParse parses the provided JSON object for a new Bounds.
func NewBoundsByParse(params map[string]interface{}) (*Bounds, error) {
	// get left, right, bottom, top extrema
	left, ok := jsonUtil.GetFloat(params, "left")
	if !ok {
		return nil, fmt.Errorf("`left` parameter missing")
	}
	right, ok := jsonUtil.GetFloat(params, "right")
	if !ok {
		return nil, fmt.Errorf("`right` parameter missing")
	}
	bottom, ok := jsonUtil.GetFloat(params, "bottom")
	if !ok {
		return nil, fmt.Errorf("`bottom` parameter missing")
	}
	top, ok := jsonUtil.GetFloat(params, "top")
	if !ok {
		return nil, fmt.Errorf("`top` parameter missing")
	}

	return NewBounds(left, right, bottom, top), nil
}

// Corners returns these extrema in the form of a binning.Bounds object
func (b Bounds) Corners() *Rectangle {
	return &Rectangle{
		BottomLeft: &Coord{
			X: b.left,
			Y: b.bottom,
		},
		TopRight: &Coord{
			X: b.right,
			Y: b.top,
		},
	}
}

// BottomLeft returns a Coord corner
func (b Bounds) BottomLeft() *Coord {
	return NewCoord(b.left, b.bottom)
}

// TopRight returns a Coord corner
func (b Bounds) TopRight() *Coord {
	return NewCoord(b.right, b.top)
}

// Left returns the left extremum
func (b Bounds) Left() float64 {
	return b.left
}

// Right returns the right extremum
func (b Bounds) Right() float64 {
	return b.right
}

// Bottom returns the bottom extremum
func (b Bounds) Bottom() float64 {
	return b.bottom
}

// Top returns the top extremum
func (b Bounds) Top() float64 {
	return b.top
}
