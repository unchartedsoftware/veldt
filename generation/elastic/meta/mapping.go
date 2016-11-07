package meta

import (
	"gopkg.in/olivere/elastic.v3"
)

// GetMapping returns the mapping for a particular elasticsearch index.
func GetMapping(client *elastic.Client, index string) (map[string]interface{}, error) {
	// get mapping
	result, err := client.
		GetMapping().
		Index(index).
		Do()
	if err != nil {
		return nil, err
	}
	return result, nil
}
