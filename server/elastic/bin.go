package elastic

/*
"aggregations": {
    "x": {
    	"buckets": [
        	{
	        	"y": {
		            "buckets": [
		            	{
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

type Row struct {
	Count uint64 `json:"doc_count"`
}

type YAgg struct {
	Rows [256]Row `json:"buckets"`
}

type Column struct {
	Y YAgg `json:"y"`
}

type XAgg struct {
	Columns [256]Column `json:"buckets"`
}

type Aggregation struct {
	X XAgg `json:"x"`
}

type Payload struct {
	Aggs Aggregation `json:"aggregations"`
}
