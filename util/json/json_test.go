package json_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"

	"github.com/unchartedsoftware/veldt/util/json"
)

var _ = Describe("json", func() {

	Describe("Get", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"obj": {}
					}
				}`)
			_, ok := json.Get(j, "test", "obj")
			Expect(ok).To(Equal(true))
		})
		It("should return the root object if no path is provided", func() {
			j := JSON(
				`{
					"test": {
						"obj": {}
					}
				}`)
			val, ok := json.Get(j)
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(j))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.Get(j, "missing", "path")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("Exists", func() {
		It("should return true if the value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"obj": {}
					}
				}`)
			exists := json.Exists(j, "test", "obj")
			Expect(exists).To(Equal(true))
		})
		It("should return false if value does not exist in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"obj": {}
					}
				}`)
			exists := json.Exists(j, "test", "missing")
			Expect(exists).To(Equal(false))
		})
	})

	Describe("GetFloat", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"float": 123.0
					}
				}`)
			_, ok := json.GetFloat(j, "test", "float")
			Expect(ok).To(Equal(true))
		})
		It("should return a float64 if the value is a float", func() {
			j := JSON(
				`{
					"test": {
						"float": 123.0
					}
				}`)
			val, ok := json.GetFloat(j, "test", "float")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(123.0))
		})
		It("should return a float64 if the value is an int", func() {
			j := JSON(
				`{
					"test": {
						"int": 123
					}
				}`)
			val, ok := json.GetFloat(j, "test", "int")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(123.0))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetFloat(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not a float", func() {
			j := JSON(
				`{
					"test": {
						"string": "hello"
					}
				}`)
			_, ok := json.GetFloat(j, "test", "string")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetInt", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"int": 123
					}
				}`)
			_, ok := json.GetInt(j, "test", "int")
			Expect(ok).To(Equal(true))
		})
		It("should return an int64 if the value is an int", func() {
			j := JSON(
				`{
					"test": {
						"int": 123
					}
				}`)
			val, ok := json.GetInt(j, "test", "int")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(int64(123)))
		})
		It("should return an int64 if the value is a float", func() {
			j := JSON(
				`{
					"test": {
						"float": 123.0
					}
				}`)
			val, ok := json.GetInt(j, "test", "float")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(int64(123)))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetInt(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an int", func() {
			j := JSON(
				`{
					"test": {
						"string": "hello"
					}
				}`)
			_, ok := json.GetInt(j, "test", "string")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetString", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"string": "hello"
					}
				}`)
			_, ok := json.GetString(j, "test", "string")
			Expect(ok).To(Equal(true))
		})
		It("should return a string if the value is a string", func() {
			j := JSON(
				`{
					"test": {
						"string": "hello"
					}
				}`)
			val, ok := json.GetString(j, "test", "string")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal("hello"))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetString(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not a string", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetString(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetBool", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"bool": true
					}
				}`)
			_, ok := json.GetBool(j, "test", "bool")
			Expect(ok).To(Equal(true))
		})
		It("should return a bool if the value is a bool", func() {
			j := JSON(
				`{
					"test": {
						"bool": false
					}
				}`)
			val, ok := json.GetBool(j, "test", "bool")
			Expect(ok).To(Equal(true))
			Expect(val).To(Equal(false))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetBool(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not a bool", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetBool(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetChild", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"child": {}
					}
				}`)
			_, ok := json.GetChild(j, "test", "child")
			Expect(ok).To(Equal(true))
		})
		It("should return a map[string]interface{} if the value is a map[string]interface{}", func() {
			j := JSON(
				`{
					"test": {
						"child": {
							"a": "a",
							"b": "b"
						}
					}
				}`)
			val, ok := json.GetChild(j, "test", "child")
			Expect(ok).To(Equal(true))
			Expect(val["a"].(string)).To(Equal("a"))
			Expect(val["b"].(string)).To(Equal("b"))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetChild(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not a map[string]interface{}", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetChild(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetArray", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"array": [0, 1, "hello", true]
					}
				}`)
			_, ok := json.GetArray(j, "test", "array")
			Expect(ok).To(Equal(true))
		})
		It("should return a []interface{} if the value is an array", func() {
			j := JSON(
				`{
					"test": {
						"array": [0, 1, "hello", true]
					}
				}`)
			val, ok := json.GetArray(j, "test", "array")
			Expect(ok).To(Equal(true))
			Expect(len(val)).To(Equal(4))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetArray(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an array", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetArray(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetFloatArray", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"array": [0, 1, 0.1, 0.2]
					}
				}`)
			_, ok := json.GetFloatArray(j, "test", "array")
			Expect(ok).To(Equal(true))
		})
		It("should return a []float64 if the value is a float array", func() {
			j := JSON(
				`{
					"test": {
						"array": [0, 1, 0.1, 0.2]
					}
				}`)
			val, ok := json.GetFloatArray(j, "test", "array")
			Expect(ok).To(Equal(true))
			Expect(val[0]).To(Equal(0.0))
			Expect(val[1]).To(Equal(1.0))
			Expect(val[2]).To(Equal(0.1))
			Expect(val[3]).To(Equal(0.2))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetFloatArray(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an array", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetFloatArray(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetIntArray", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"array": [0, 1, 0.1, 0.2]
					}
				}`)
			_, ok := json.GetIntArray(j, "test", "array")
			Expect(ok).To(Equal(true))
		})
		It("should return a []int64 if the value is an int array", func() {
			j := JSON(
				`{
					"test": {
						"array": [0, 1, 0.1, 0.2]
					}
				}`)
			val, ok := json.GetIntArray(j, "test", "array")
			Expect(ok).To(Equal(true))
			Expect(val[0]).To(Equal(int64(0)))
			Expect(val[1]).To(Equal(int64(1)))
			Expect(val[2]).To(Equal(int64(0)))
			Expect(val[3]).To(Equal(int64(0)))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetIntArray(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an array", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetIntArray(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetStringArray", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"array": ["a", "b", "see", "dee"]
					}
				}`)
			_, ok := json.GetStringArray(j, "test", "array")
			Expect(ok).To(Equal(true))
		})
		It("should return a []string if the value is a string array", func() {
			j := JSON(
				`{
					"test": {
						"array": ["a", "b", "see", "dee"]
					}
				}`)
			val, ok := json.GetStringArray(j, "test", "array")
			Expect(ok).To(Equal(true))
			Expect(val[0]).To(Equal("a"))
			Expect(val[1]).To(Equal("b"))
			Expect(val[2]).To(Equal("see"))
			Expect(val[3]).To(Equal("dee"))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetStringArray(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an array", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetStringArray(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetBoolArray", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"array": [true, false, false, true]
					}
				}`)
			_, ok := json.GetBoolArray(j, "test", "array")
			Expect(ok).To(Equal(true))
		})
		It("should return a []bool if the value is a bool array", func() {
			j := JSON(
				`{
					"test": {
						"array": [true, false, false, true]
					}
				}`)
			val, ok := json.GetBoolArray(j, "test", "array")
			Expect(ok).To(Equal(true))
			Expect(val[0]).To(Equal(true))
			Expect(val[1]).To(Equal(false))
			Expect(val[2]).To(Equal(false))
			Expect(val[3]).To(Equal(true))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetBoolArray(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an array", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetBoolArray(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetChildArray", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"array": [{}, {}]
					}
				}`)
			_, ok := json.GetChildArray(j, "test", "array")
			Expect(ok).To(Equal(true))
		})
		It("should return a []map[string]interface{} if the value is an array of nodes", func() {
			j := JSON(
				`{
					"test": {
						"array": [{
							"a": "a"
						}]
					}
				}`)
			val, ok := json.GetChildArray(j, "test", "array")
			Expect(ok).To(Equal(true))
			Expect(val[0]["a"]).To(Equal("a"))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetChildArray(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an array", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetChildArray(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetChildMap", func() {
		It("should return a true value if a value exists in the provided path", func() {
			j := JSON(
				`{
					"test": {
						"children": {
							"a": {},
							"b": {},
							"c": {}
						}
					}
				}`)
			_, ok := json.GetChildMap(j, "test", "children")
			Expect(ok).To(Equal(true))
		})
		It("should return a map[string]map[string]interface{} if the value is a map of nodes", func() {
			j := JSON(
				`{
					"test": {
						"children": {
							"a": {
								"val": "a"
							},
							"b": {
								"val": "b"
							},
							"c": {
								"val": "c"
							}
						}
					}
				}`)
			val, ok := json.GetChildMap(j, "test", "children")
			Expect(ok).To(Equal(true))
			Expect(val["a"]["val"]).To(Equal("a"))
			Expect(val["b"]["val"]).To(Equal("b"))
			Expect(val["c"]["val"]).To(Equal("c"))
		})
		It("should return a false value if value does not exist in the provided path", func() {
			j := JSON(`{}`)
			_, ok := json.GetChildMap(j, "test", "missing")
			Expect(ok).To(Equal(false))
		})
		It("should return a false value if value is not an array", func() {
			j := JSON(
				`{
					"test": {
						"int": 5
					}
				}`)
			_, ok := json.GetChildMap(j, "test", "int")
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetRandomChild", func() {
		It("should return a true if there is at least one nested object", func() {
			j := JSON(
				`{
					"test": {
						"a": {},
						"b": {},
						"c": {}
					}
				}`)
			_, ok := json.GetRandomChild(j, "test")
			Expect(ok).To(Equal(true))
		})
		It("should return a map[string]interface{} if there is at least one nested object", func() {
			j := JSON(
				`{
					"test": {
						"child" : {
							"a": "a"
						}
					}
				}`)
			val, ok := json.GetRandomChild(j, "test")
			Expect(ok).To(Equal(true))
			Expect(val["a"]).To(Equal("a"))
		})
		It("should return a false if there are no nested objects", func() {
			j := JSON(
				`{
					"test": {}
				}`)
			_, ok := json.GetRandomChild(j, "test")
			Expect(ok).To(Equal(false))
		})
	})

})
