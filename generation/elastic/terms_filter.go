package elastic

import (
	"strings"

	"gopkg.in/olivere/elastic.v3"
)

// GetTermsFilter creates and returns a simple terms filter for elastic based
// tiling.
func GetTermsFilter(arg map[string]interface{}) (interface{}, bool) {
	field, fieldOk := arg["field"].(string)
	terms, termsOk := arg["terms"].(string)
	if fieldOk && termsOk {
		split := strings.Split(terms, ",")
		ts := make([]interface{}, len(split))
		for _, str := range split {
			ts = append(ts, str)
		}
		return elastic.NewTermsQuery(field, ts...), true
	}
	return nil, false
}
