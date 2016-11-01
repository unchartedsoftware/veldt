package tile

import (
	"fmt"

	"github.com/unchartedsoftware/prism/binning"
)

// Request represents the tile type and tile coord.
type Request struct {
	Coord  *binning.TileCoord     `json:"coord"`
	Type   string                 `json:"type"`
	URI    string                 `json:"uri"`
	Store  string                 `json:"store"`
	Params map[string]interface{} `json:"params"`
}

// String returns the request formatted as a string.
func (r *Request) String() string {
	return fmt.Sprintf("%s/%s/%d/%d/%d",
		r.Type,
		r.URI,
		r.Coord.Z,
		r.Coord.X,
		r.Coord.Y)
}
