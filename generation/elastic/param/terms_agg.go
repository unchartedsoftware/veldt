package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultTextField = "text"
)

// TermsAgg represents params for extracting term counts.
type TermsAgg struct {
	Field string
	Terms []string
}

// NewTermsAgg instantiates and returns a new parameter object.
func NewTermsAgg(tileReq *tile.Request) (*TermsAgg, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "terms_agg")
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("TermsAgg `field` parameter missing from tiling request %s", tileReq.String())
	}
	terms, ok := json.GetStringArray(params, "terms")
	if !ok {
		return nil, fmt.Errorf("TermsAgg `terms` parameter missing from tiling request %s", tileReq.String())
	}
	return &TermsAgg{
		Field: field,
		Terms: terms,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *TermsAgg) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.Field,
		strings.Join(p.Terms, ":"))
}

// GetQuery returns an elastic query.
func (p *TermsAgg) GetQuery() *elastic.TermsQuery {
	terms := make([]interface{}, len(p.Terms))
	for i, term := range p.Terms {
		terms[i] = term
	}
	return elastic.NewTermsQuery(p.Field, terms...)
}

// GetAggregations returns an elastic aggregation.
func (p *TermsAgg) GetAggregations() map[string]*elastic.FilterAggregation {
	aggs := make(map[string]*elastic.FilterAggregation, len(p.Terms))
	// add all filter aggregations
	for _, term := range p.Terms {
		aggs[term] = elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery(p.Field, term))
	}
	return aggs
}
