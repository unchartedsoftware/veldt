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

	Describe("BufferValue", func() {
		It("should add the stringified value to the internal buffer at the appropriate indent", func() {
			validator.BufferValue("test", nil)
			Expect(validator.String()).To(Equal(`"test"`))
		})
		It("should annotate the stringified value with the provided error", func() {
			err := fmt.Errorf("error")
			validator.BufferValue("test", err)
			str := "vvvvvv\n" +
				"\"test\"\n" +
				"^^^^^^ Error: " + err.Error()
			Expect(color.RemoveColor(validator.String())).To(Equal(str))
		})
	})

	Describe("String", func() {
		It("should stringify and return the internal buffer", func() {
			validator.BufferValue("test", nil)
			Expect(validator.String()).To(Equal(`"test"`))
		})
	})

	Describe("Size", func() {
		It("should return the line count of the internal buffer", func() {
			validator.BufferValue("0", nil)
			validator.BufferValue("1", nil)
			validator.BufferValue("2", nil)
			validator.BufferValue("3", nil)
			Expect(validator.Size()).To(Equal(4))
		})
	})

	Describe("StartObject", func() {
		It("should append an open bracket (`{`) in the buffer", func() {
			validator.StartObject()
			str := "{"
			Expect(validator.String()).To(Equal(str))
		})
		It("should increment the next indentation", func() {
			validator.StartObject()
			validator.BufferValue("test", nil)
			validator.EndObject()
			str := "{\n    \"test\"\n}"
			Expect(validator.String()).To(Equal(str))
		})
	})

	Describe("EndObject", func() {
		It("should append an open bracket (`}`) in the buffer", func() {
			validator.StartObject()
			validator.EndObject()
			str := "{\n}"
			Expect(validator.String()).To(Equal(str))
		})
		It("should decrement the next indentation", func() {
			validator.StartObject()
			validator.EndObject()
			validator.StartObject()
			validator.EndObject()
			str := "{\n},\n{\n}"
			Expect(validator.String()).To(Equal(str))
		})
	})

	Describe("StartArray", func() {
		It("should append an open bracket (`{`) in the buffer", func() {
			validator.StartArray()
			str := "["
			Expect(validator.String()).To(Equal(str))
		})
		It("should increment the next indentation", func() {
			validator.StartArray()
			validator.BufferValue("test", nil)
			validator.EndArray()
			str := "[\n    \"test\"\n]"
			Expect(validator.String()).To(Equal(str))
		})
	})

	Describe("EndArray", func() {
		It("should append an open bracket (`}`) in the buffer", func() {
			validator.StartArray()
			validator.EndArray()
			str := "[\n]"
			Expect(validator.String()).To(Equal(str))
		})
		It("should decrement the next indentation", func() {
			validator.StartArray()
			validator.EndArray()
			validator.StartArray()
			validator.EndArray()
			str := "[\n],\n[\n]"
			Expect(validator.String()).To(Equal(str))
		})
	})

	Describe("BufferKeyValue", func() {
		It("should buffer a string key and a number", func() {
			str :=
				"{\n" +
					"    \"a\": 0.5\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.String()).To(Equal(str))
		})
		It("should buffer a string key and a string", func() {
			str :=
				"{\n" +
					"    \"a\": \"test\"\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.String()).To(Equal(str))
		})
		It("should buffer a string key and an array", func() {
			str :=
				"{\n" +
					"    \"a\": [\n" +
					"        0.5,\n" +
					"        \"test\"\n" +
					"    ]\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.String()).To(Equal(str))
		})
		It("should buffer a string key and an object", func() {
			str :=
				"{\n" +
					"    \"a\": {\n" +
					"        \"b\": [\n" +
					"            0.5,\n" +
					"            \"test\",\n" +
					"            true,\n" +
					"            false\n" +
					"        ]\n" +
					"    }\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.String()).To(Equal(str))
		})
		It("should handle nesting objects with proper indentation", func() {
			str :=
				"{\n" +
					"    \"a\": {\n" +
					"        \"b\": {\n" +
					"            \"c\": {\n" +
					"                \"d\": [\n" +
					"                    0.5,\n" +
					"                    \"test\"\n" +
					"                ]\n" +
					"            }\n" +
					"        }\n" +
					"    }\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.String()).To(Equal(str))
		})
		It("should handle nesting arrays with proper indentation", func() {
			str :=
				"{\n" +
					"    \"a\": [\n" +
					"        0.5,\n" +
					"        [\n" +
					"            [\n" +
					"                0.5,\n" +
					"                \"test\"\n" +
					"            ],\n" +
					"            [\n" +
					"                0.5,\n" +
					"                \"test\"\n" +
					"            ]\n" +
					"        ]\n" +
					"    ]\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.String()).To(Equal(str))
		})
		It("should handle nesting objects in arrays with proper indentation", func() {
			str :=
				"{\n" +
					"    \"a\": [\n" +
					"        {\n" +
					"            \"a\": 2\n" +
					"        },\n" +
					"        {\n" +
					"            \"a\": 2\n" +
					"        }\n" +
					"    ]\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.String()).To(Equal(str))
		})
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
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], err)
			validator.EndObject()
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
			validator.StartObject()
			validator.StartError(err.Error())
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndError()
			validator.EndObject()
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
			validator.StartObject()
			validator.StartError(err)
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndError()
			validator.EndObject()
			Expect(validator.HasError()).To(Equal(true))
		})
		It("should return false if the validator did not processes any errors", func() {
			str :=
				"{\n" +
					"    \"a\": \"???\"\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.HasError()).To(Equal(false))
		})
	})

	Describe("Error", func() {
		It("should return an error with the annoted JSON as the string", func() {
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
			validator.StartObject()
			validator.StartError(err.Error())
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndError()
			validator.EndObject()
			err = validator.Error()
			Expect(err).NotTo(BeNil())
			Expect(color.RemoveColor(err.Error())).To(Equal(output))
		})
		It("should return nil if the validator did not processes any errors", func() {
			str :=
				"{\n" +
					"    \"a\": \"???\"\n" +
					"}"
			j := JSON(str)
			validator.StartObject()
			validator.BufferKeyValue("a", j["a"], nil)
			validator.EndObject()
			Expect(validator.Error()).To(BeNil())
		})
	})

})
