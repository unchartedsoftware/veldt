package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

type boolQuery struct {
}

// BoolQuery represents params for a boolean query on a tile
type BoolQuery struct {
	Query *elastic.BoolQuery
}

// NewBoolQuery instantiates and returns a new parameter object.
func NewBoolQuery(tileReq *tile.Request) (*BoolQuery, error) {
	param, ok := json.GetChild(tileReq.Params, "bool_query")
	if !ok {
		fmt.Printf("bool_query parameter missing from tiling request %s\n", tileReq.String())
		return nil, fmt.Errorf("BoolQuery parameter missing from tiling request %s", tileReq.String())
	}

	musts, ok := json.GetChildrenArray(param, "must")
	if !ok {
		fmt.Println("couldn't find musts")
	}
	// allocate a new list of term queries the size of musts
	mustQueries := make([]elastic.Query, len(musts))

	for i, must := range musts {
		termQueryDef, ok := json.GetChild(must, "term")
		if ok {
			field, _ := json.GetString(termQueryDef, "field")
			termsList, _ := json.GetStringArray(termQueryDef, "terms")
			mustQueries[i] = elastic.NewTermsQuery(field, termsList)
			continue
		}
		rangeQueryDef, ok := json.GetChild(must, "range")
		if ok {
			field, _ := json.GetString(rangeQueryDef, "field")
			from, _ := json.GetNumber(rangeQueryDef, "from") // really only need one of 'from' or 'to'
			to, _ := json.GetNumber(rangeQueryDef, "to")
			mustQueries[i] = elastic.NewRangeQuery(field).From(from).To(to)
			continue
		}
	}

	bq := elastic.NewBoolQuery()

	for _, query := range mustQueries {
		bq.Must(query)
	}

	// TODO add support for the other parts of the bool query
	// shoulds, ok := json.GetChildrenArray(param, "shoulds")
	// mustNots, ok := json.GetChildrenArray(param, "must_nots")
	// filters, ok := json.GetChildrenArray(param, "filters")

	return &BoolQuery{
		Query: bq,
	}, nil
}

// GetHash will return a hash of params
func (bq *BoolQuery) GetHash() string {
	q, _ := bq.Query.Source()
	// fmt.Println(q)
	top := q.(map[string]interface{})

	var hashes []string

	// get a list of the nodes available under must
	nodes, found := json.GetChildrenArray(top, "bool", "must")

	if found {
		for _, node := range nodes {
			term, ok := json.GetChild(node, "terms")
			if ok {
				for field, list := range term {
					actualTerms := list.([]interface{})[0]
					hash := fmt.Sprintf("%s:%v", field, actualTerms)
					hashes = append(hashes, hash)
				}
			}
			rangeJSON, ok := json.GetChild(node, "range")
			if ok {
				for key, val := range rangeJSON {
					rangeBlob := val.(map[string]interface{})
					hash := fmt.Sprintf("%s:%v:%v", key, rangeBlob["from"], rangeBlob["to"])
					hashes = append(hashes, hash)
				}
			}
		}
		return strings.Join(hashes, ":")
	}

	fmt.Println("Unable to get hash of bool query params")
	return "not-a-hash"
}
