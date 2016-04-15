package param_test

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/tile"
)

var _ = Describe("bool_query", func() {

	const (
		jsonParams1 = `{
		    "binning":{
		        "x":"locality_bag.dateBegin",
		        "left":1350416916775,
		        "right":1446058904529,
		        "y":"feature.author_firstmessage_rank",
		        "bottom":0,
		        "top":78891,
		        "resolution":256
		    },
		    "bool_query":{
		        "must":[
		            {
		                "term":{
		                    "field":"feature.firearm_type",
		                    "terms":[
		                        "Rifle",
		                        "Handgun"
		                    ]
		                }
		            },
		            {
		                "range":{
		                    "field":"feature.indicator_risky",
		                    "from":0.5,
		                    "to":3
		                }
		            }
		        ],
		        "must_not":[

		        ],
		        "should":[

		        ],
		        "filter":[

		        ]
		    }
		}`
		jsonParams2 = `{
		    "binning":{
		        "x":"locality_bag.dateBegin",
		        "left":1350416916775,
		        "right":1446058904529,
		        "y":"feature.author_firstmessage_rank",
		        "bottom":0,
		        "top":78891,
		        "resolution":256
		    },
		    "bool_query":{
		        "must":[
		            {
		                "term":{
		                    "field":"feature.author",
		                    "terms":[
		                        "Penn",
		                        "Teller"
		                    ]
		                }
		            }
		        ]
		    }
		}`
	)

	var (
		tileReq tile.Request
		bq      *param.BoolQuery
	)

	Describe("GetHash", func() {
		It("Should hash correctly for multiple queries in must clause", func() {
			decoder := json.NewDecoder(strings.NewReader(jsonParams1))
			tileReq.Params = make(map[string]interface{})
			err := decoder.Decode(&tileReq.Params)
			if err != nil {
				fmt.Println("error decoding json")
			}
			bq, _ = param.NewBoolQuery(&tileReq)
			hash := bq.GetHash()
			Expect(hash).To(Equal("feature.firearm_type:Rifle:Handgun::feature.indicator_risky:0.500000:3.000000"))
		})

		It("Should hash correctly for single query in must clause", func() {
			decoder := json.NewDecoder(strings.NewReader(jsonParams2))
			tileReq.Params = make(map[string]interface{})
			err := decoder.Decode(&tileReq.Params)
			if err != nil {
				fmt.Println("error decoding json")
			}
			bq, _ = param.NewBoolQuery(&tileReq)
			hash := bq.GetHash()
			Expect(hash).To(Equal("feature.author:Penn:Teller"))
		})
	})
})
