package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tiling"
)

// Extrema represents the min and max values for an ordinal property.
type Extrema struct {
	Min float64
	Max float64
}

// GetExtrema returns the extrema of a numeric field for the provided index.
func GetExtrema(endpoint string, index string, field string) (*Extrema, error) {
	// get client
	client, err := getClient(endpoint)
	if err != nil {
		return nil, err
	}
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
		return nil, fmt.Errorf("Min '%s' aggregation was not found in response", field )
	}
	max, ok := result.Aggregations.Max("max")
	if !ok {
		return nil, fmt.Errorf("Max '%s' aggregation was not found in response", field )
	}
	return &Extrema{
		Min: *min.Value,
		Max: *max.Value,
	}, nil
}
