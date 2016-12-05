package citus

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/tile"
)

type TopHits struct {
	tile.TopHits
}

func (t *TopHits) AddAggs(query *Query) *Query {
	//Select the top N rows when sorted. Return only the specified fields.
	for _, field := range t.IncludeFields {
		query.Select(field)
	}

	// sort
	if t.SortOrder == "desc" {
		query.OrderBy(fmt.Sprintf("%s DESC", t.SortField))
	} else {
		query.OrderBy(t.SortField)
	}

	query.Limit(uint32(t.HitsCount))

	return query
}

func (t *TopHits) GetTopHits(rows *pgx.Rows) ([]map[string]interface{}, error) {
	hits := make([]map[string]interface{}, 0)
	for rows.Next() {
		columnValues, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rowResult := make(map[string]interface{})

		//Cycle through the fields to create the map.
		for i, field := range t.IncludeFields {
			rowResult[field] = columnValues[i]
		}

		hits = append(hits, rowResult)
	}
	return hits, nil
}
