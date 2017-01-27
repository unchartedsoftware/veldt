package citus

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/veldt/tile"
)

// TargetTerms represents a citus implementation of the target terms tile.
type TargetTerms struct {
	tile.TargetTerms
}

// AddQuery adds the tiling query to the provided query object.
func (t *TargetTerms) AddQuery(query *Query) *Query {
	//Want to keep only documents that have the specified terms.
	//Use the already existing Has construct.
	hasQuery := &Has{}
	hasQuery.Field = t.TermsField
	terms := make([]interface{}, len(t.Terms))
	for i, term := range t.Terms {
		terms[i] = term
	}
	hasQuery.Values = terms

	clause, _ := hasQuery.Get(query)
	query.Where(clause)
	return query
}

// AddAggs adds the tiling aggregations to the provided query object.
func (t *TargetTerms) AddAggs(query *Query) *Query {
	//Count by term, only considering the specified terms.
	//Assume the backing field is an array. Need to unpack that array and group by the terms.
	query.Select(fmt.Sprintf("unnest(%s) AS term", t.TermsField))

	query.GroupBy("term")
	query.Select("COUNT(*) as term_count")
	query.OrderBy("term_count desc")

	//Generate the filter for the terms.
	clause := ""
	for _, value := range t.Terms {
		valueParam := query.AddParameter(value)
		clause = clause + fmt.Sprintf(", %s", valueParam)
	}
	query.Where(fmt.Sprintf("term IN [%s]", clause[2:]))

	return query
}

// GetTerms parses the result of the terms query into a map of term -> count.
func (t *TargetTerms) GetTerms(rows *pgx.Rows) (map[string]uint32, error) {
	// build map of topics and counts
	counts := make(map[string]uint32)
	for rows.Next() {
		var term string
		var count uint32
		err := rows.Scan(&term, &count)
		if err != nil {
			return nil, fmt.Errorf("Error parsing top terms: %v", err)
		}
		counts[term] = count
	}
	return counts, nil
}
