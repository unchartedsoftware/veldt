package json_test

import (
	j "encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unchartedsoftware/veldt/util/json"
)

var _ = Describe("json", func() {

	const (
		data = `{
			"nested": {
				"child": {
					"number": 2323423
				}
			},
			"binning": {
				"x":"locality_bag.dateBegin",
				"left":1350416916775,
				"right":1446058904529,
				"y":"feature.author_firstmessage_rank",
				"bottom":0,
				"top":78891,
				"resolution":256
			},
			"bool_query": {
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
		jsonNode   map[string]interface{}
		emptyNode  map[string]interface{}
		nested     map[string]interface{}
		binning    map[string]interface{}
		bool_query map[string]interface{}
	)

	BeforeEach(func() {
		// empty node
		jsonNode = make(map[string]interface{})
		// json node
		jsonNode = make(map[string]interface{})
		err := j.Unmarshal([]byte(data), &jsonNode)
		if err != nil {
			fmt.Println("Error decoding json")
		}
		Expect(err).To(BeNil())
		// extract nodes
		v := jsonNode["nested"]
		nested = v.(map[string]interface{})
		v = jsonNode["binning"]
		binning = v.(map[string]interface{})
		v = jsonNode["bool_query"]
		bool_query = v.(map[string]interface{})
	})

	Describe("Get", func() {
		It("Should return a true value if a value exists in the provided path", func() {
			val, ok := json.Get(jsonNode, "binning", "left")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(binning["left"]))
		})
		It("Should return the root object if no path is provided", func() {
			val, ok := json.Get(jsonNode)
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(jsonNode))
		})
		It("Should return a false value if value does not exist in the provided path", func() {
			_, ok := json.Get(jsonNode, "binning", "left", "incorrect")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("Exists", func() {
		It("Should return a true value if the value exists in the provided path", func() {
			exists := json.Exists(jsonNode, "binning", "left")
			Expect(exists).To(Equal(true))
		})
		It("Should return a false value if value does not exist in the provided path", func() {
			exists := json.Exists(jsonNode, "binning", "left", "incorrect")
			Expect(exists).To(Equal(false))
		})
	})

	Describe("GetChild", func() {
		It("Should return a true value if an object value exists in the provided path", func() {
			val, ok := json.GetChild(jsonNode, "nested", "child")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(nested["child"]))
		})
		It("Should return the root object if no path is provided", func() {
			val, ok := json.GetChild(jsonNode)
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(jsonNode))
		})
		It("Should return a false value if an object value does not exist in the provided path or is not an object", func() {
			_, ok := json.GetChild(jsonNode, "binning", "left")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetRandomChild", func() {
		It("Should return a true if there is at least one nested object", func() {
			val, ok := json.GetRandomChild(jsonNode)
			Expect(ok).To(Equal(true))
			Expect(val).To(Or(Equal(nested), Equal(binning), Equal(bool_query)))
		})
		It("Should return a false if there are no nested objects", func() {
			_, ok := json.GetRandomChild(emptyNode)
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetChildOrEmpty", func() {
		It("Should return the nested object if it exists", func() {
			v := json.GetChildOrEmpty(jsonNode, "binning")
			Expect(v).To(Equal(binning))
		})
		It("Should return an empty map if the node does not exist", func() {
			v := json.GetChildOrEmpty(emptyNode, "child")
			Expect(v).To(Not(BeNil()))
		})
		It("Should return the root object if no path is provided", func() {
			val := json.GetChildOrEmpty(jsonNode)
			Expect(val).To(Equal(jsonNode))
		})
	})

	Describe("GetChildArray", func() {
		It("Should return a list of nodes", func() {
			jsonChildArray, ok := json.GetChildArray(jsonNode, "bool_query", "must")
			Expect(ok).To(Equal(true))
			Expect(len(jsonChildArray)).To(Equal(2))
		})
		It("should work with only a single path element", func() {
			jsonChild, ok := json.GetChild(jsonNode, "bool_query")
			Expect(ok).To(Equal(true))
			jsonChildArray, ok := json.GetChildArray(jsonChild, "must")
			Expect(ok).To(Equal(true))
			Expect(len(jsonChildArray)).To(Equal(2))
		})
	})
})
