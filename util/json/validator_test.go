package json_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"

	"github.com/unchartedsoftware/veldt/util/color"
	"github.com/unchartedsoftware/veldt/util/json"
)

var _ = Describe("validator", func() {

	var validator *json.Validator

	BeforeEach(func() {
		validator = &json.Validator{}
	})

	Describe("Buffer", func() {
		It("should add the string to the internal buffer at the provided input", func() {
			validator.Buffer("test", 0)
			Expect(validator.String()).To(Equal("test"))
		})
		It("should indent the buffered string by the magnitude of the second parameter", func() {
			validator.Buffer("test", 2)
			Expect(validator.String()).To(Equal("        test"))
		})
	})

	Describe("String", func() {
		It("should stringify and return the results", func() {
			validator.Buffer("test", 0)
			Expect(validator.String()).To(Equal("test"))
		})
	})

	Describe("Size", func() {
		It("should return the line count of the internal buffer", func() {
			validator.Buffer("0", 0)
			validator.Buffer("1", 0)
			validator.Buffer("2", 0)
			validator.Buffer("3", 0)
			Expect(validator.Size()).To(Equal(4))
		})
	})

	Describe("BufferKeyValue", func() {
		It("should buffer a string key and a number", func() {
			str :=
				"{\n" +
					"    \"a\": 0.5\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.Buffer("}", 0)
			Expect(validator.String()).To(Equal(str))
		})
		It("should buffer a string key and a string", func() {
			str :=
				"{\n" +
					"    \"a\": \"test\"\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.Buffer("}", 0)
			Expect(validator.String()).To(Equal(str))
		})
		It("should buffer a string key and an array", func() {
			str :=
				"{\n" +
					"    \"a\": [ 0.5, \"test\" ]\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.Buffer("}", 0)
			Expect(validator.String()).To(Equal(str))
		})
		It("should buffer a string key and an object", func() {
			str :=
				"{\n" +
					"    \"a\": {\n" +
					"        \"b\": [ 0.5, \"test\" ]\n" +
					"    }\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.Buffer("}", 0)
			Expect(validator.String()).To(Equal(str))
		})
		It("should handle nesting objects with proper indentation", func() {
			str :=
				"{\n" +
					"    \"a\": {\n" +
					"        \"b\": {\n" +
					"            \"c\": {\n" +
					"                \"d\": [ 0.5, \"test\" ]\n" +
					"            }\n" +
					"        }\n" +
					"    }\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.Buffer("}", 0)
			Expect(validator.String()).To(Equal(str))
		})
		It("should handle nesting arrays with proper indentation", func() {
			str :=
				"{\n" +
					"    \"a\": [ 0.5, [ [ 0.5, \"test\" ], [ 0.5, \"test\" ] ] ]\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.Buffer("}", 0)
			Expect(validator.String()).To(Equal(str))
		})
		// It("should handle nesting objects in arrays with proper indentation", func() {
		// 	str :=
		// 		"{\n" +
		// 		"    \"a\": [\n" +
		// 		"        {\n" +
		// 		"            \"a\": 2\n" +
		// 		"        },\n" +
		// 		"        {\n" +
		// 		"            \"a\": 2\n" +
		// 		"        }\n" +
		// 		"    ]\n" +
		// 		"}"
		// 	j := JSON(str)
		// 	validator.Buffer("{", 0)
		// 	validator.BufferKeyValue("a", j["a"], 1, nil)
		// 	validator.Buffer("}", 0)
		// 	fmt.Println()
		// 	fmt.Println(validator.String())
		// 	Expect(validator.String()).To(Equal(str))
		// })
		It("should annotate the key and value with an optional error", func() {
			str :=
				"{\n" +
					"    \"a\": \"???\"\n" +
					"}"
			err := fmt.Errorf("missing value")
			output :=
				"{\n" +
					"    vvvvvvvvvv\n" +
					"    \"a\": \"???\"\n" +
					"    ^^^^^^^^^^ Error: " + err.Error() + "\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.BufferKeyValue("a", j["a"], 1, err)
			validator.Buffer("}", 0)
			Expect(color.RemoveColor(validator.String())).To(Equal(output))
		})
	})

	Describe("StartError / EndError", func() {
		It("should wrap the intermediate lines with an error annotation", func() {
			str :=
				"{\n" +
					"    \"a\": \"???\"\n" +
					"}"
			err := fmt.Errorf("missing value")
			output :=
				"{\n" +
					"    vvvvvvvvvv\n" +
					"    \"a\": \"???\"\n" +
					"    ^^^^^^^^^^ Error: " + err.Error() + "\n" +
					"}"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.StartError(err.Error(), 1)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.EndError()
			validator.Buffer("}", 0)
			Expect(color.RemoveColor(validator.String())).To(Equal(output))
		})
	})

	Describe("HasError", func() {
		It("should return true if the validator processes any errors", func() {
			str :=
				"{\n" +
					"    \"a\": \"???\"\n" +
					"}"
			err := "missing value"
			j := JSON(str)
			validator.Buffer("{", 0)
			validator.StartError(err, 1)
			validator.BufferKeyValue("a", j["a"], 1, nil)
			validator.EndError()
			validator.Buffer("}", 0)
		})
	})

})
