package citus

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/veldt/tile"
)

// TopTerms represents a citus implementation of the top terms tile.
type TopTerms struct {
	tile.TopTerms
}

// AddAggs adds the tiling aggregations to the provided query object.
func (t *TopTerms) AddAggs(query *Query) *Query {
	//Assume the backing field is an array. Need to unpack that array and group by the terms.
	query.Select(fmt.Sprintf("unnest(%s) AS term", t.TermsField))

	query.GroupBy("term")
	query.Select("COUNT(*) as term_count")
	query.OrderBy("term_count desc")
	query.Limit(uint32(t.TermsCount))

	return query
}

// GetTerms parses the result of the terms query into a map of term -> count.
func (t *TopTerms) GetTerms(rows *pgx.Rows) (map[string]uint32, error) {
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
