package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

type prefixQuery struct {
	Field    string
	Prefixes []string
}

func (q *prefixQuery) getQueries() []*elastic.PrefixQuery {
	prefixes := make([]*elastic.PrefixQuery, len(q.Prefixes))
	for i, prefix := range q.Prefixes {
		prefixes[i] = elastic.NewPrefixQuery(q.Field, prefix)
	}
	return prefixes
}

func (q *prefixQuery) getHash() string {
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(q.Prefixes, ":"))
}

// PrefixFilter represents params for filtering by certain prefixes.
type PrefixFilter struct {
	Queries []*prefixQuery
}

// NewPrefixFilter instantiates and returns a new parameter object.
func NewPrefixFilter(tileReq *tile.Request) (*PrefixFilter, error) {
	params, ok := json.GetChildrenArray(tileReq.Params, "prefix_filter")
	if !ok {
		return nil, ErrMissing
	}
	// parse each range query
	queries := make([]*prefixQuery, len(params))
	for i, param := range params {
		field, ok := json.GetString(param, "field")
		if !ok {
			return nil, fmt.Errorf("PrefixFilter `field` parameter missing from tiling request %s", tileReq.String())
		}
		prefixes, ok := json.GetStringArray(param, "prefixes")
		if !ok {
			return nil, fmt.Errorf("PrefixFilter `prefixes` parameter missing from tiling request %s", tileReq.String())
		}
		queries[i] = &prefixQuery{
			Field:    field,
			Prefixes: prefixes,
		}
	}
	return &PrefixFilter{
		Queries: queries,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *PrefixFilter) GetHash() string {
	strs := make([]string, len(p.Queries))
	for i, query := range p.Queries {
		strs[i] = query.getHash()
	}
	return strings.Join(strs, "::")
}

// GetQueries returns a slice of elastic queries.
func (p *PrefixFilter) GetQueries() []*elastic.PrefixQuery {
	var queries []*elastic.PrefixQuery
	for _, query := range p.Queries {
		queries = append(queries, query.getQueries()...)
	}
	return queries
}
