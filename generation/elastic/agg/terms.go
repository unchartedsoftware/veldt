package agg

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultSize = 10
)

// Terms represents a terms aggregation.
type Terms struct {
	Field string
	Size  int
}

// NewTerms instantiates and returns a new parameter object.
func NewTerms(params map[string]interface{}) (*Terms, error) {
	params, ok := json.GetChild(params, "terms")
	if !ok {
		return nil, fmt.Errorf("%s `terms` aggregation parameter", param.MissingPrefix)
	}
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Terms `field` parameter missing from tiling param %v", params)
	}
	size := int(json.GetNumberDefault(params, defaultSize, "size"))
	return &Terms{
		Field: field,
		Size:  size,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Terms) GetHash() string {
	return fmt.Sprintf("%s:%d",
		p.Field,
		p.Size,
	)
}

// GetAgg returns an elastic aggregation.
func (p *Terms) GetAgg() *elastic.TermsAggregation {
	return elastic.NewTermsAggregation().
		Field(p.Field).
		Size(p.Size)
}
