package agg

import (
	"fmt"

	"github.com/unchartedsoftware/prism/generation/citus/param"
	"github.com/unchartedsoftware/prism/generation/citus/query"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultTermsSize = 10
)

// TopTerms represents params for extracting particular topics.
type TopTerms struct {
	Field string
	Size  uint32
}

// NewTopTerms instantiates and returns a new topic parameter object.
func NewTopTerms(params map[string]interface{}) (*TopTerms, error) {
	params, ok := json.GetChild(params, "top_terms")
	if !ok {
		return nil, fmt.Errorf("%s `top_terms` aggregation parameter", param.MissingPrefix)
	}
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("TopTerms `field` parameter missing from tiling param %v", params)
	}
	return &TopTerms{
		Field: field,
		Size:  uint32(json.GetNumberDefault(params, defaultTermsSize, "size")),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *TopTerms) GetHash() string {
	return fmt.Sprintf("%s:%d", p.Field, p.Size)
}

// AddAgg returns an elastic aggregation.
func (p *TopTerms) AddAgg(q *query.Query) (*query.Query, error) {
	//Add the terms field. Assume the field is an array.
	q.AddField(fmt.Sprintf("unnest(%s) AS term", p.Field))

	//Need to nest the existing query as a table and group by the terms.
	termQuery, err := query.NewQuery()
	if err != nil {
		return nil, err
	}

	termQuery.AddTable(fmt.Sprintf("(%s) terms", q.GetQuery(true)))
	termQuery.AddGroupByClause("term")
	termQuery.AddField("term")
	termQuery.AddField("COUNT(*) as term_count")
	termQuery.AddOrderByClause("term_count desc")
	termQuery.SetLimit(p.Size)

	return termQuery, nil
}
