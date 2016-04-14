package param

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultTermsSize = 10
)

// TopTerms represents params for extracting particular topics.
type TopTerms struct {
	Field string
	Size  uint32
}

// NewTopTerms instantiates and returns a new topic parameter object.
func NewTopTerms(tileReq *tile.Request) (*TopTerms, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "top_terms")
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("TopTerms `field` parameter missing from tiling request %s", tileReq.String())
	}
	return &TopTerms{
		Field: field,
		Size:  uint32(json.GetNumberDefault(params, "size", defaultTermsSize)),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *TopTerms) GetHash() string {
	return fmt.Sprintf("%s:%d", p.Field, p.Size)
}

// GetAggregation returns an elastic aggregation.
func (p *TopTerms) GetAggregation() *elastic.TermsAggregation {
	return elastic.NewTermsAggregation().
		Field(p.Field).
		Size(int(p.Size))
}
