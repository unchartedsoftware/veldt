package agg

import (
	"fmt"

	"github.com/unchartedsoftware/prism/util/json"
)

// CustomAggs represents params for a custom es aggregation.
type CustomAggs struct {
	Aggs map[string]interface{}
}

// NewCustomAggs instantiates and returns a new topic parameter object.
func NewCustomAggs(params map[string]interface{}) (*CustomAggs, error) {
	params, ok := json.GetChild(params, "custom_aggs")
	if !ok {
		return nil, fmt.Errorf("%v `custom_aggs` aggregation parameter", params)
	}
	aggs, ok := json.GetChild(params, "aggs")
	if !ok {
		return nil, fmt.Errorf("aggs was not of type `map[string]interface{}` in response for request %v",
			params)
	}
	return &CustomAggs{
		Aggs: aggs,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *CustomAggs) GetHash() string {
	return json.GetHash(p.Aggs)
}

// GetAgg returns an elastic aggregation.
func (p *CustomAggs) GetAgg() map[string]interface{} {
	return p.Aggs
}
