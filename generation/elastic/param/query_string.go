package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

type queryStringQuery struct {
	Field  string
	String string
}

func (q *queryStringQuery) getQuery() elastic.Query {
	return elastic.NewQueryStringQuery(q.String).
		Field(q.Field)
}

func (q *queryStringQuery) getHash() string {
	return fmt.Sprintf("%s:%s",
		q.Field,
		q.String)
}

// QueryString represents params for filtering by a query string.
type QueryString struct {
	Queries []*queryStringQuery
}

// NewQueryString instantiates and returns a new query string parameter object.
func NewQueryString(tileReq *tile.Request) (*QueryString, error) {
	params, ok := json.GetChildrenArray(tileReq.Params, "query_string")
	if !ok {
		return nil, fmt.Errorf("QueryString parameter missing from tiling request %s", tileReq.String())
	}
	// parse each range query
	queries := make([]*queryStringQuery, len(params))
	for i, param := range params {
		field, ok := json.GetString(param, "field")
		if !ok {
			return nil, fmt.Errorf("QueryString `field` parameter missing from tiling request %s", tileReq.String())
		}
		str, ok := json.GetString(param, "string")
		if !ok {
			return nil, fmt.Errorf("QueryString `string` parameter missing from tiling request %s", tileReq.String())
		}
		queries[i] = &queryStringQuery{
			Field:  field,
			String: str,
		}
	}
	return &QueryString{
		Queries: queries,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *QueryString) GetHash() string {
	strs := make([]string, len(p.Queries))
	for i, query := range p.Queries {
		strs[i] = query.getHash()
	}
	return strings.Join(strs, "::")
}

// GetQueries returns a slice of elastic queries.
func (p *QueryString) GetQueries() []elastic.Query {
	queries := make([]elastic.Query, len(p.Queries))
	for i, query := range p.Queries {
		queries[i] = query.getQuery()
	}
	return queries
}
