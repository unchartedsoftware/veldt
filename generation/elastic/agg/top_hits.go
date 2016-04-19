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
		return nil, param.ErrMissing
	}
	size := int(json.GetNumberDefault(params, "size", defaultHitsSize))
	srt := json.GetStringDefault(params, "sort", "")
	order := json.GetStringDefault(params, "order", defaultOrder)
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
	if len(p.Sort) > 0 {
		if p.Order == "desc" {
			agg.Sort(p.Sort, false)
		} else {
			agg.Sort(p.Sort, true)
		}
	}
	// ...
	return agg
}
