package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

type rangeQuery struct {
	Field string
	From  interface{}
	To    interface{}
}

func (q *rangeQuery) getQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(q.Field).
		Gte(q.From).
		Lte(q.To)
}

func (q *rangeQuery) getHash() string {
	return fmt.Sprintf("%s:%f:%f",
		q.Field,
		q.From,
		q.To)
}

// Range represents params for restricting the range of a given field in a tile.
type Range struct {
	Queries []*rangeQuery
}

// NewRange instantiates and returns a new range parameter object.
func NewRange(tileReq *tile.Request) (*Range, error) {
	params, ok := json.GetChildrenArray(tileReq.Params, "range")
	if !ok {
		return nil, ErrMissing
	}
	// parse each range query
	queries := make([]*rangeQuery, len(params))
	for i, param := range params {
		field, ok := json.GetString(param, "field")
		if !ok {
			return nil, fmt.Errorf("Range `field` parameter missing from tiling request %s", tileReq.String())
		}
		from, ok := json.GetInterface(param, "from")
		if !ok {
			return nil, fmt.Errorf("Range `from` parameter missing from tiling request %s", tileReq.String())
		}
		to, ok := json.GetInterface(param, "to")
		if !ok {
			return nil, fmt.Errorf("Range `to` parameter missing from tiling request %s", tileReq.String())
		}
		queries[i] = &rangeQuery{
			Field: field,
			From:  from,
			To:    to,
		}
	}
	return &Range{
		Queries: queries,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Range) GetHash() string {
	strs := make([]string, len(p.Queries))
	for i, query := range p.Queries {
		strs[i] = query.getHash()
	}
	return strings.Join(strs, "::")
}

// GetQueries returns a slice of elastic queries.
func (p *Range) GetQueries() []*elastic.RangeQuery {
	queries := make([]*elastic.RangeQuery, len(p.Queries))
	for i, query := range p.Queries {
		queries[i] = query.getQuery()
	}
	return queries
}
