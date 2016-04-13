package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

type termQuery struct {
	Field string
	Terms []string
}

func (q *termQuery) getQuery() elastic.Query {
	terms := make([]interface{}, len(q.Terms))
	for i, term := range q.Terms {
		terms[i] = term
	}
	return elastic.NewTermsQuery(q.Field, terms...)
}

func (q *termQuery) getHash() string {
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(q.Terms, ":"))
}

// TermsFilter represents params for filtering by certain terms.
type TermsFilter struct {
	Queries []*termQuery
}

// NewTermsFilter instantiates and returns a new parameter object.
func NewTermsFilter(tileReq *tile.Request) (*TermsFilter, error) {
	params, ok := json.GetChildrenArray(tileReq.Params, "terms_filter")
	if !ok {
		return nil, fmt.Errorf("TermsFilter parameter missing from tiling request %s", tileReq.String())
	}
	// parse each range query
	queries := make([]*termQuery, len(params))
	for i, param := range params {
		field, ok := json.GetString(param, "field")
		if !ok {
			return nil, fmt.Errorf("TermsFilter `field` parameter missing from tiling request %s", tileReq.String())
		}
		terms, ok := json.GetStringArray(param, "terms")
		if !ok {
			return nil, fmt.Errorf("TermsFilter `terms` parameter missing from tiling request %s", tileReq.String())
		}
		queries[i] = &termQuery{
			Field: field,
			Terms: terms,
		}
	}
	return &TermsFilter{
		Queries: queries,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *TermsFilter) GetHash() string {
	strs := make([]string, len(p.Queries))
	for i, query := range p.Queries {
		strs[i] = query.getHash()
	}
	return strings.Join(strs, "::")
}

// GetQueries returns a slice of elastic queries.
func (p *TermsFilter) GetQueries() []elastic.Query {
	queries := make([]elastic.Query, len(p.Queries))
	for i, query := range p.Queries {
		queries[i] = query.getQuery()
	}
	return queries
}
