package citus

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/tile"
)

// TopTerms represents a tiling generator that produces heatmaps.
type TopTerms struct {
	tile.TopTerms
}

func (t *TopTerms) AddAggs(query *Query) *Query {
	//TODO: Find a better way to make this work. The caller NEEDS to use the returned value.
	//Assume the backing field is an array. Need to unpack that array and group by the terms.
	query.Select(fmt.Sprintf("unnest(%s) AS term", t.TermsField))

	//Need to nest the existing query as a table and group by the terms.
	//TODO: Figure out how to handle error properly.
	termQuery, _ := NewQuery()

	termQuery.From(fmt.Sprintf("(%s) terms", query.GetQuery(true)))
	termQuery.GroupBy("term")
	termQuery.Select("term")
	termQuery.Select("COUNT(*) as term_count")
	termQuery.OrderBy("term_count desc")
	termQuery.Limit(uint32(t.TermsCount))

	return termQuery
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
