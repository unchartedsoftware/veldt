package rest

import (
	"github.com/unchartedsoftware/prism/tile"
)

// TileGenerator represents a base generator that uses elasticsearch for its
// backend.
type TileGenerator struct {
	req *tile.Request
}

// GetHash returns the hash for this generator.
func (g *TileGenerator) GetHash() string {
	return ""
}
