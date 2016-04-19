package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
)

// GetExtrema returns the extrema of a numeric field for the provided index.
func GetExtrema(client *elastic.Client, index string, field string) (*binning.Extrema, error) {
	// query
	result, err := client.
		Search(index).
		Size(0).
		Aggregation("min",
		elastic.NewMinAggregation().
			Field(field)).
		Aggregation("max",
		elastic.NewMaxAggregation().
			Field(field)).
		Do()
	if err != nil {
		return nil, err
	}
	// parse aggregations
	min, ok := result.Aggregations.Min("min")
	if !ok {
		return nil, fmt.Errorf("Min '%s' aggregation was not found in response for %s", field, index)
	}
	max, ok := result.Aggregations.Max("max")
	if !ok {
		return nil, fmt.Errorf("Max '%s' aggregation was not found in response for %s", field, index)
	}
	// it seems if the mapping exists, but no documents have the attribute, the min / max are null
	if min.Value == nil || max.Value == nil {
		return nil, nil
	}
	return &binning.Extrema{
		Min: *min.Value,
		Max: *max.Value,
	}, nil
}
