package citus

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/generation/citus/agg"
	"github.com/unchartedsoftware/prism/generation/citus/param"
	"github.com/unchartedsoftware/prism/generation/citus/query"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	termsAggName     = "topterms"
	histogramAggName = "histogramAgg"
)

// TopCountTile represents a tiling generator that produces top term counts.
type TopCountTile struct {
	TileGenerator
	Tiling   *param.Tiling
	TopTerms *agg.TopTerms
	Query    *query.Query
}

// NewTopCountTile instantiates and returns a pointer to a new generator.
func NewTopCountTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		citus, err := param.NewCitus(tileReq)
		if err != nil {
			return nil, err
		}
		// required
		tiling, err := param.NewTiling(tileReq)
		if err != nil {
			return nil, err
		}
		topTerms, err := agg.NewTopTerms(tileReq.Params)
		if err != nil {
			return nil, err
		}
		query, err := query.NewQuery()
		if err != nil {
			return nil, err
		}
		t := &TopCountTile{}
		t.Citus = citus
		t.Tiling = tiling
		t.TopTerms = topTerms
		t.Query = query
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopCountTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.TopTerms,
		g.Query,
	}
}

func (g *TopCountTile) getQuery(q *query.Query) *query.Query {
	g.Tiling.AddXQuery(q)
	g.Tiling.AddYQuery(q)
	return q
}

func (g *TopCountTile) getAgg(q *query.Query) (*query.Query, error) {
	// get top terms agg
	aq, err := g.TopTerms.AddAgg(q)
	if err != nil {
		return nil, err
	}
	return aq, nil
}

func (g *TopCountTile) parseResult(rows *pgx.Rows) ([]byte, error) {
	// Result of query is term, count.
	counts := make(map[string]interface{})
	for rows.Next() {
		var term string
		var count float64
		err := rows.Scan(&term, &count)
		if err != nil {
			return nil, fmt.Errorf("Error parsing top terms: %s %v",
				g.req.String(), err)
		}
		counts[term] = count
	}
	// marshal results map
	return json.Marshal(counts)
}

// GetTile returns the marshalled tile data.
func (g *TopCountTile) GetTile() ([]byte, error) {
	// send query
	query := g.Query
	query.AddTable(g.req.URI)
	query = g.getQuery(query)
	query, err := g.getAgg(query)
	rows, err := g.client.Query(query.GetQuery(false), query.QueryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// parse and return results
	return g.parseResult(rows)
}
