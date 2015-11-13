package elastic

/*
"hits": {
	"total": 5793,
	"max_score": 1,
	"hits": [
		{
			"_index": "isil_twitter_dec2may",
			"_type": "datum",
			"_id": "tweet561584461964668928",
			"_score": 1,
			"fields": {
				"locality: {
					location": "40.759831,-73.988277"
				}
			}
		}
		...
	]
}
*/

type JsonLocation struct {
	Location string `json:"location"`
}

type JsonFields struct {
	Locality JsonLocation `json:"locality"`
}

type JsonBin struct {
	Source JsonFields `json:"_source"`
}

type JsonBinSet struct {
	Bins []JsonBin `json:"hits"`
}

type JsonPayload struct {
	Hits JsonBinSet `json:"hits"`
}
