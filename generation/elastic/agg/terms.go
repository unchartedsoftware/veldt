package agg

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultTextField = "text"
)

// Terms represents params for extracting term counts.
type Terms struct {
	Field string
	Terms []string
}

// NewTerms instantiates and returns a new parameter object.
func NewTerms(params map[string]interface{}) (*Terms, error) {
	params, ok := json.GetChild(params, "terms")
	if !ok {
		return nil, param.ErrMissing
	}
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Terms `field` parameter missing from tiling param %v", params)
	}
	terms, ok := json.GetStringArray(params, "terms")
	if !ok {
		return nil, fmt.Errorf("Terms `terms` parameter missing from tiling param %v", params)
	}
	sort.Strings(terms)
	return &Terms{
		Field: field,
		Terms: terms,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Terms) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.Field,
		strings.Join(p.Terms, ":"))
}

// GetQuery returns an elastic query.
func (p *Terms) GetQuery() *elastic.TermsQuery {
	terms := make([]interface{}, len(p.Terms))
	for i, term := range p.Terms {
		terms[i] = term
	}
	return elastic.NewTermsQuery(p.Field, terms...)
}

// GetAggregations returns an elastic aggregation.
func (p *Terms) GetAggregations() map[string]*elastic.FilterAggregation {
	aggs := make(map[string]*elastic.FilterAggregation, len(p.Terms))
	// add all filter aggregations
	for _, term := range p.Terms {
		aggs[term] = elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery(p.Field, term))
	}
	return aggs
}
