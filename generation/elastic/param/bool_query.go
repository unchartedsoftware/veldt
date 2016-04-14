package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// BoolQuery represents params for a boolean query on a tile
type BoolQuery struct {
	query *elastic.BoolQuery
	hash  string
}

type queryComponent interface {
	getQuery() elastic.Query
	getHash() string
}

// NewBoolQuery instantiates and returns a new parameter object.
func NewBoolQuery(tileReq *tile.Request) (*BoolQuery, error) {
	musts, ok := json.GetChildrenArray(tileReq.Params, "bool_query", "must")
	if !ok {
		return nil, fmt.Errorf("bool_query/must clause path not found in request params %s", tileReq.String())
	}
	// allocate a new list of term queries the size of musts
	mustQueryComponents := make([]queryComponent, len(musts))

	for i, must := range musts {
		termQueryDef, ok := json.GetChild(must, "term")
		if ok {
			field, _ := json.GetString(termQueryDef, "field")
			termsList, _ := json.GetStringArray(termQueryDef, "terms")
			mustQueryComponents[i] = &termQuery{field, termsList}
			continue
		}
		rangeQueryDef, ok := json.GetChild(must, "range")
		if ok {
			field, _ := json.GetString(rangeQueryDef, "field")
			from, _ := json.GetNumber(rangeQueryDef, "from") // TODO really only need one of 'from' or 'to'
			to, _ := json.GetNumber(rangeQueryDef, "to")
			mustQueryComponents[i] = &rangeQuery{field, from, to}
			continue
		}
	}

	bq := elastic.NewBoolQuery()

	var hashes []string

	for _, query := range mustQueryComponents {
		hashes = append(hashes, query.getHash())
		bq.Must(query.getQuery())
	}

	hash := strings.Join(hashes, "::")

	// TODO add support for the other parts of the bool query
	// shoulds, ok := json.GetChildrenArray(param, "shoulds")
	// mustNots, ok := json.GetChildrenArray(param, "must_nots")
	// filters, ok := json.GetChildrenArray(param, "filters")

	return &BoolQuery{
		query: bq,
		hash:  hash,
	}, nil
}

// GetHash will return the calculated hash of bool query params
func (bq *BoolQuery) GetHash() string {
	return bq.hash
}

// GetQuery will return the elastic query object
func (bq *BoolQuery) GetQuery() elastic.Query {
	return bq.query
}
