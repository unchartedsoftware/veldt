package geometry

import (
	"fmt"
	"math"

	jsonUtil "github.com/unchartedsoftware/veldt/util/json"
)

// Bounds represents a bounding rectangle with convenience getters.
type Bounds struct {
	Left   float64
	Right  float64
	Bottom float64
	Top    float64
}

// NewBounds returns a new Bounds object.
func NewBounds(left, right, bottom, top float64) *Bounds {
	return &Bounds{
		Left:   left,
		Right:  right,
		Bottom: bottom,
		Top:    top,
	}
}

// Parse parses the provided JSON object and populates the bounds attributes.
func (b *Bounds) Parse(params map[string]interface{}) error {
	// get left, right, bottom, top extrema
	left, ok := jsonUtil.GetFloat(params, "left")
	if !ok {
		return fmt.Errorf("`left` parameter missing")
	}
	right, ok := jsonUtil.GetFloat(params, "right")
	if !ok {
		return fmt.Errorf("`right` parameter missing")
	}
	bottom, ok := jsonUtil.GetFloat(params, "bottom")
	if !ok {
		return fmt.Errorf("`bottom` parameter missing")
	}
	top, ok := jsonUtil.GetFloat(params, "top")
	if !ok {
		return fmt.Errorf("`top` parameter missing")
	}
	b.Left = left
	b.Right = right
	b.Bottom = bottom
	b.Top = top
	return nil
}

// MinX returns the minimum x value for the bounds.
func (b Bounds) MinX() float64 {
	return math.Min(b.Left, b.Right)
}

// MaxX returns the maximum x value for the bounds.
func (b Bounds) MaxX() float64 {
	return math.Max(b.Left, b.Right)
}

// MinY returns the minimum y value for the bounds.
func (b Bounds) MinY() float64 {
	return math.Min(b.Bottom, b.Top)
}

// MaxY returns the maximum y value for the bounds.
func (b Bounds) MaxY() float64 {
	return math.Max(b.Bottom, b.Top)
}

// RangeX returns the absolute distance between left and right.
func (b *Bounds) RangeX() float64 {
	return math.Abs(b.Right - b.Left)
}

// RangeY returns the absolute distance between top and bottom.
func (b *Bounds) RangeY() float64 {
	return math.Abs(b.Top - b.Bottom)
}
