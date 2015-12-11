package elastic

// GetMapping returns the mapping for a particular elasticsearch index.
func GetMapping(endpoint string, index string) (map[string]interface{}, error) {
	// get client
	client, err := getClient(endpoint)
	if err != nil {
		return nil, err
	}
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
