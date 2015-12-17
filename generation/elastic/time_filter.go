package elastic

import (
	"gopkg.in/olivere/elastic.v3"
)

// GetTimeFilter creates and returns a simple time filter for elastic based
// tiling.
func GetTimeFilter(arg map[string]interface{}) (interface{}, bool) {
	field, fieldOk := arg["field"].(string)
	from, fromOk := arg["from"].(string)
	to, toOk := arg["to"].(string)
	if fieldOk && fromOk && toOk {
		return elastic.NewRangeQuery(field).
			Gte(from).
			Lte(to), true
	}
	return nil, false
}
