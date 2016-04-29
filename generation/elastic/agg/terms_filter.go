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

// TermsFilter represents a terms filter aggregation.
type TermsFilter struct {
	Field string
	Terms []string
}

// NewTermsFilter instantiates and returns a new parameter object.
func NewTermsFilter(params map[string]interface{}) (*TermsFilter, error) {
	params, ok := json.GetChild(params, "terms_filter")
	if !ok {
		return nil, fmt.Errorf("%s `terms_filter` aggregation parameter", param.MissingPrefix)
	}
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("TermsFilter `field` parameter missing from tiling param %v", params)
	}
	terms, ok := json.GetStringArray(params, "terms")
	if !ok {
		return nil, fmt.Errorf("TermsFilter `terms` parameter missing from tiling param %v", params)
	}
	sort.Strings(terms)
	return &TermsFilter{
		Field: field,
		Terms: terms,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *TermsFilter) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.Field,
		strings.Join(p.Terms, ":"))
}

// GetQuery returns an elastic query.
func (p *TermsFilter) GetQuery() *elastic.TermsQuery {
	terms := make([]interface{}, len(p.Terms))
	for i, term := range p.Terms {
		terms[i] = term
	}
	return elastic.NewTermsQuery(p.Field, terms...)
}

// GetAggs returns the elastic aggregations.
func (p *TermsFilter) GetAggs() map[string]*elastic.FilterAggregation {
	aggs := make(map[string]*elastic.FilterAggregation, len(p.Terms))
	// add all filter aggregations
	for _, term := range p.Terms {
		aggs[term] = elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery(p.Field, term))
	}
	return aggs
}
