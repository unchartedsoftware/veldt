package param

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// Elastic represents params specifically tailored to elasticsearch queries.
type Elastic struct {
	Types []string
}

// NewElastic instantiates and returns a new macro/micro parameter object.
func NewElastic(tileReq *tile.Request) (*Elastic, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "elastic")
	types, _ := json.GetStringArray(params, "types")
	sort.Strings(types)
	return &Elastic{
		Types: types,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Elastic) GetHash() string {
	return fmt.Sprintf("%s",
		strings.Join(p.Types, ":"))
}

// GetSearchService returns an elastic.SearchService based on the params.
func (p *Elastic) GetSearchService(client *elastic.Client) *elastic.SearchService {
	return client.Search().Type(p.Types...)
}
