package citus

func EncodeFrequency(frequency []*FrequencyResult) []map[string]interface{} {
	buckets := make([]map[string]interface{}, len(frequency))
	for i, bucket := range frequency {
		buckets[i] = map[string]interface{}{
			"timestamp": bucket.Bucket,
			"count":     bucket.Value,
		}
	}

	return buckets
}
