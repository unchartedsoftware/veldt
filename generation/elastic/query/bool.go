package query

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Bool represents a boolean query.
type Bool struct {
	musts    []Query
	mustNots []Query
	shoulds  []Query
}

// Query represents a base query Query interface.
type Query interface {
	GetQuery() elastic.Query
	GetHash() string
}

func getQueryByType(query map[string]interface{}) (Query, error) {
	// bool
	params, ok := json.GetChild(query, "bool")
	if ok {
		return NewBool(params)
	}
	// exists
	params, ok = json.GetChild(query , "exists")
	if (ok) {
		return NewExists(params)
	}
	// terms
	params, ok = json.GetChild(query, "terms")
	if ok {
		return NewTerms(params)
	}
	// range
	params, ok = json.GetChild(query, "range")
	if ok {
		return NewRange(params)
	}
	// prefix
	params, ok = json.GetChild(query, "prefix")
	if ok {
		return NewPrefix(params)
	}
	// query_string
	params, ok = json.GetChild(query, "query_string")
	if ok {
		return NewString(params)
	}
	return nil, fmt.Errorf("No recognized query type found in %v", query)
}

// NewBool instantiates and returns a new parameter object.
func NewBool(params map[string]interface{}) (*Bool, error) {
	// must queries
	must, ok := json.GetChildArray(params, "must")
	var musts []Query
	if ok {
		musts = make([]Query, len(must))
		for i, query := range must {
			q, err := getQueryByType(query)
			if err != nil {
				return nil, err
			}
			musts[i] = q
		}
	} else {
		musts = make([]Query, 0)
	}
	// must not queries
	mustNot, ok := json.GetChildArray(params, "must_not")
	var mustNots []Query
	if ok {
		mustNots = make([]Query, len(mustNot))
		for i, query := range mustNot {
			q, err := getQueryByType(query)
			if err != nil {
				return nil, err
			}
			mustNots[i] = q
		}
	} else {
		mustNots = make([]Query, 0)
	}
	// should queries
	should, ok := json.GetChildArray(params, "should")
	var shoulds []Query
	if ok {
		shoulds = make([]Query, len(should))
		for i, query := range should {
			q, err := getQueryByType(query)
			if err != nil {
				return nil, err
			}
			shoulds[i] = q
		}
	} else {
		shoulds = make([]Query, 0)
	}
	return &Bool{
		musts:    musts,
		mustNots: mustNots,
		shoulds:  shoulds,
	}, nil
}

// GetHash will return the hash of the parameters.
func (b *Bool) GetHash() string {
	numMusts := len(b.musts)
	numMustNots := len(b.mustNots)
	numShoulds := len(b.shoulds)
	hashes := make([]string, numMusts+numMustNots+numShoulds)
	for i, query := range b.musts {
		hashes[i] = query.GetHash()
	}
	for i, query := range b.mustNots {
		hashes[i+numMusts] = query.GetHash()
	}
	for i, query := range b.shoulds {
		hashes[i+numMusts+numMustNots] = query.GetHash()
	}
	return strings.Join(hashes, "::")
}

// GetQuery will return the elastic query object.
func (b *Bool) GetQuery() elastic.Query {
	query := elastic.NewBoolQuery()
	// add musts
	for _, must := range b.musts {
		query.Must(must.GetQuery())
	}
	// add must nots
	for _, mustNot := range b.mustNots {
		query.MustNot(mustNot.GetQuery())
	}
	// add shoulds
	for _, should := range b.shoulds {
		query.Should(should.GetQuery())
	}
	return query
}
