package json_test

import (
	j "encoding/json"
	"fmt"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unchartedsoftware/prism/util/json"
)

var _ = Describe("json", func() {

	const (
		jsonRequest = `{
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
	            "must": [
	                {
	                    "term": {
	                        "field": "feature.firearm_type",
	                        "terms":  ["Rifle"]
	                    }
	                },
	                {
	                    "range": {
	                        "field": "feature.indicator_risky",
	                        "from": 0.5,
	                        "to": 3
	                    }
	                }
	            ],
	            "must_not": [],
	            "should": [],
	            "filter": []
	        }
	    }`
	)

	var (
		jsonNode json.Node
	)

	BeforeEach(func() {
		decoder := j.NewDecoder(strings.NewReader(jsonRequest))
		jsonNode = make(map[string]interface{})
		err := decoder.Decode(&jsonNode)
		if err == io.EOF {
			fmt.Println("Error decoding json")
		}
	})

	Describe("GetChildrenArray", func() {
		It("Should return a list of nodes", func() {
			jsonChildArray, ok := json.GetChildrenArray(jsonNode, "bool_query", "must")
			Expect(ok).To(Equal(true))
			Expect(len(jsonChildArray)).To(Equal(2))
		})
		It("should work with only a single path element", func() {
			jsonChild, ok := json.GetChild(jsonNode, "bool_query")
			Expect(ok).To(Equal(true))
			jsonChildArray, ok := json.GetChildrenArray(jsonChild, "must")
			Expect(ok).To(Equal(true))
			Expect(len(jsonChildArray)).To(Equal(2))
		})
	})
})
