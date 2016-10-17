package citus

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/generation/citus/param"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// TileGenerator represents a base generator that uses citus for its backend.
type TileGenerator struct {
	host   string
	port   string
	client *pgx.ConnPool
	req    *tile.Request
	Citus  *param.Citus
}

// GetHash returns the hash for this generator.
func (g *TileGenerator) GetHash() string {
	return fmt.Sprintf("%s:%s:%s",
		g.host,
		g.port,
		g.Citus.GetHash())
}
