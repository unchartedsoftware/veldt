package citus

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/meta"
)

// MetaGenerator represents a base generator that uses elasticsearch for its
// backend.
type MetaGenerator struct {
	host   string
	port   string
	client *pgx.ConnPool
	req    *meta.Request
}

// GetHash returns the hash for this generator.
func (g *MetaGenerator) GetHash() string {
	return fmt.Sprintf("%s:%s", g.host, g.port)
}
