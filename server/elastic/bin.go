package elastic

/*
"aggregations": {
    "x": {
    	"buckets": [
        	{
				"key": 1261961216,
	        	"y": {
		            "buckets": [
		            	{
							"key": 1615331328,"
			                "doc_count": 10
			            },
						...
					]
				}
			},
			...
		]
	}
}
*/

type Topic struct {
	Count uint64 `json:"doc_count"`
}

type TopicPayload struct {
	Aggs map[string]Topic `json:"aggregations"`
}

///

type Row struct {
	Count uint64 `json:"doc_count"`
	PixelY uint64 `json:"key"`
}

type YAgg struct {
	Rows []Row `json:"buckets"`
}

type Column struct {
	Y YAgg `json:"y"`
	PixelX uint64 `json:"key"`
}

type XAgg struct {
	Columns []Column `json:"buckets"`
}

type Aggregation struct {
	X XAgg `json:"x"`
}

type HeatmapPayload struct {
	Aggs Aggregation `json:"aggregations"`
}
