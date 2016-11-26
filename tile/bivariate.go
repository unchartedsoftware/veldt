package tile

import (
	"fmt"

	"github.com/unchartedsoftware/prism/util/json"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	XField     string
	YField     string
	Left       float64
	Right      float64
	Bottom     float64
	Top        float64
	Resolution int
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (b *Bivariate) Parse(params map[string]interface{}) error {
	// get x and y fields
	xField, ok := json.GetString(params, "xField")
	if !ok {
		return fmt.Errorf("`xField` parameter missing from tile")
	}
	yField, ok := json.GetString(params, "yField")
	if !ok {
		return fmt.Errorf("`yField` parameter missing from tile")
	}
	// get left, right, bottom, top extrema
	left, ok := json.GetNumber(params, "left")
	if !ok {
		return fmt.Errorf("`left` parameter missing from tile")
	}
	right, ok := json.GetNumber(params, "right")
	if !ok {
		return fmt.Errorf("`right` parameter missing from tile")
	}
	bottom, ok := json.GetNumber(params, "bottom")
	if !ok {
		return fmt.Errorf("`bottom` parameter missing from tile")
	}
	top, ok := json.GetNumber(params, "top")
	if !ok {
		return fmt.Errorf("`top` parameter missing from tile")
	}
	// get resolution
	resolution := json.GetNumberDefault(params, 256, "resolution")
	// set attributes
	b.XField = xField
	b.YField = yField
	b.Left = left
	b.Right = right
	b.Bottom = bottom
	b.Top = top
	b.Resolution = int(resolution)
	return nil
}
