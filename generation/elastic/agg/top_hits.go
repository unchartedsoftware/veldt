package agg

import (
	ejson "encoding/json"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultHitsSize = 25
	defaultOrder    = "desc"
)

// TopHits represents params for binning the data within the tile.
type TopHits struct {
	Size    int
	Include []string
	Sort    string
	Order   string
}

// NewTopHits instantiates and returns a new metric aggregation parameter.
func NewTopHits(params map[string]interface{}) (*TopHits, error) {
	params, ok := json.GetChild(params, "top_hits")
	if !ok {
		return nil, fmt.Errorf("%s `top_hits` aggregation parameter", param.MissingPrefix)
	}
	size := int(json.GetNumberDefault(params, defaultHitsSize, "size"))
	srt := json.GetStringDefault(params, "", "sort")
	order := json.GetStringDefault(params, defaultOrder, "order")
	include, ok := json.GetStringArray(params, "include")
	if !ok {
		include = nil
	} else {
		sort.Strings(include)
	}
	return &TopHits{
		Size:    size,
		Include: include,
		Sort:    srt,
		Order:   order,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *TopHits) GetHash() string {
	return fmt.Sprintf("%d:%s:%s:%s",
		p.Size,
		p.Sort,
		p.Order,
		strings.Join(p.Include, ":"))
}

// GetAgg returns an elastic aggregation.
func (p *TopHits) GetAgg() elastic.Aggregation {
	agg := elastic.NewTopHitsAggregation().
		Size(p.Size)
	// sort
	if len(p.Sort) > 0 {
		if p.Order == "desc" {
			agg.Sort(p.Sort, false)
		} else {
			agg.Sort(p.Sort, true)
		}
	}
	// add includes
	if p.Include != nil {
		agg.FetchSourceContext(
			elastic.NewFetchSourceContext(true).
				Include(p.Include...))
	}
	return agg
}

// GetHitsMap parses and unmarshals the top hits.
func (p *TopHits) GetHitsMap(agg *elastic.AggregationTopHitsMetric) ([]map[string]interface{}, bool) {
	topHits := make([]map[string]interface{}, len(agg.Hits.Hits))
	for index, hit := range agg.Hits.Hits {
		if err := ejson.Unmarshal(*hit.Source, &(topHits[index])); err != nil {
			return nil, false
		}
	}
	return topHits, true
}
