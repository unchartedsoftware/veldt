package citus

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/veldt/tile"
)

// TermsFrequency represents a citus implementation of the terms frequency tile.
type TermsFrequency struct {
	tile.TermsFrequency
}

// AddAggs adds the tiling aggregations to the provided query object.
func (t *TermsFrequency) AddAggs(query *Query) *Query {

	//Count by term
	if t.FieldType == "string" {
		query.Select(fmt.Sprintf("%s AS term", t.TermsField))
		query.Where(fmt.Sprintf("%s IS NOT NULL", t.TermsField))
	} else {
		//Assume the backing field is an array. Need to unpack that array and group by the terms.
		query.Select(fmt.Sprintf("unnest(%s) AS term", t.TermsField))
	}

	query.GroupBy("term")
	query.Select("COUNT(*) as term_count")
	query.OrderBy("term_count desc")

	return query
}

// GetTerms parses the result of the terms query into a map of term -> count.
func (t *TermsFrequency) GetTerms(rows *pgx.Rows) (map[string]uint32, error) {
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
