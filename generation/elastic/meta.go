package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/meta"
)

// MetaGenerator represents a base generator that uses elasticsearch for its
// backend.
type MetaGenerator struct {
	host   string
	port   string
	client *elastic.Client
	req    *meta.Request
}

// GetHash returns the hash for this generator.
func (g *MetaGenerator) GetHash() string {
	return fmt.Sprintf("%s:%s", g.host, g.port)
}
