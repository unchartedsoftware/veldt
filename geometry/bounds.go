package geometry

import (
	"fmt"
	"math"

	jsonUtil "github.com/unchartedsoftware/veldt/util/json"
)

// Bounds represents an Immutable bounding rectangle with convenience getters.
type Bounds struct {
	left               float64
	right              float64
	bottom             float64
	top                float64
	minX               float64
	maxX               float64
	minY               float64
	maxY               float64
	rangeX             float64
	rangeY             float64
	isMinMaxCalculated bool
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

func (b *Bounds) calculateMinMax() {
	b.minX = math.Min(b.left, b.right)
	b.maxX = math.Max(b.left, b.right)
	b.minY = math.Min(b.bottom, b.top)
	b.maxY = math.Max(b.bottom, b.top)
	b.rangeX = math.Abs(b.right - b.left)
	b.rangeY = math.Abs(b.top - b.bottom)
	b.isMinMaxCalculated = true
}

// MinX returns the minimum value of Left and Right
func (b Bounds) MinX() float64 {
	if !b.isMinMaxCalculated {
		b.calculateMinMax()
	}
	return b.minX
}

// MaxX returns the maximum value of Left and Right
func (b Bounds) MaxX() float64 {
	if !b.isMinMaxCalculated {
		b.calculateMinMax()
	}
	return b.maxX
}

// MinY returns the minimum value of Bottom and Top
func (b Bounds) MinY() float64 {
	if !b.isMinMaxCalculated {
		b.calculateMinMax()
	}
	return b.minY
}

// MaxY returns the maximum value of Bottom and Top
func (b Bounds) MaxY() float64 {
	if !b.isMinMaxCalculated {
		b.calculateMinMax()
	}
	return b.maxY
}

// RangeX returns the absolute distance between left and right
func (b Bounds) RangeX() float64 {
	if !b.isMinMaxCalculated {
		b.calculateMinMax()
	}
	return b.rangeX
}

// RangeY returns the absolute distance between bottom and top
func (b Bounds) RangeY() float64 {
	if !b.isMinMaxCalculated {
		b.calculateMinMax()
	}
	return b.rangeY
}
