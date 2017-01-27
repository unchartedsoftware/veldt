package tile

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// Frequency represents a tile which returns data over a provided time range.
type Frequency struct {
	FrequencyField string
	GT             interface{}
	GTE            interface{}
	LT             interface{}
	LTE            interface{}
	Interval       string
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *Frequency) Parse(params map[string]interface{}) error {
	frequencyField, ok := json.GetString(params, "frequencyField")
	if !ok {
		return fmt.Errorf("`frequencyField` parameter missing from tile")
	}
	gte, gteOk := json.Get(params, "gte")
	gt, gtOk := json.Get(params, "gt")
	lte, lteOk := json.Get(params, "lte")
	lt, ltOk := json.Get(params, "lt")
	if !gteOk && !gtOk && !lteOk && !ltOk {
		return fmt.Errorf("top term frequency has no valid range parameters")
	}
	if gteOk && gtOk {
		return fmt.Errorf("both `gte` and `gt` have been provided, only one upper bound may be provided")
	}
	if lteOk && ltOk {
		return fmt.Errorf("both `lte` and `lt` have been provided, only one lower bound may be provided")
	}
	interval, ok := json.GetString(params, "interval")
	if !ok {
		return fmt.Errorf("`interval` parameter missing from tile")
	}
	t.FrequencyField = frequencyField
	t.GTE = gte
	t.GT = gt
	t.LTE = lte
	t.LT = lt
	t.Interval = interval
	return nil
}
