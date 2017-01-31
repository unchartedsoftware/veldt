package tile

import (
	"fmt"

	"github.com/unchartedsoftware/prism/util/json"
)

// BivariatePair represents the parameters required for any bivariate-pair tile.
type BivariatePair struct {
	Bivariate
	X2Field    string
	Y2Field    string
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (b *BivariatePair) Parse(params map[string]interface{}) error {
	b.Bivariate.Parse(params)

	x2Field, ok := json.GetString(params, "x2Field")
	if !ok {
		return fmt.Errorf("`x2Field` parameter missing from tile")
	}
	y2Field, ok := json.GetString(params, "y2Field")
	if !ok {
		return fmt.Errorf("`y2Field` parameter missing from tile")
	}
	b.X2Field = x2Field
	b.Y2Field = y2Field
	fmt.Printf("<><><> tile.bivariate_pair: Parse completed. b.Right %f \n", b.Right)
	return nil
}

// GetX2Y2 given a data hit, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *BivariatePair) GetX2Y2(hit map[string]interface{}) (float32, float32, bool) { // TODO: DRY this out.
	// get x / y fields from data
	ix, ok := hit[b.X2Field]
	if !ok {
		return 0, 0, false
	}
	iy, ok := hit[b.Y2Field]
	if !ok {
		return 0, 0, false
	}
	// get X / Y of the data
	x, y, ok := castPixel(ix, iy)
	if !ok {
		return 0, 0, false
	}
	// convert to tile pixel coords in the range [0 - 256)
	tx := b.Bivariate.GetX(x)
	ty := b.Bivariate.GetY(y)
	// return position in tile coords with 2 decimal places
	return toFixed(float32(tx), 2), toFixed(float32(ty), 2), true
}

