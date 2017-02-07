package test

import (
	"encoding/json"
	"github.com/onsi/gomega"
)

// JSON converts the provided string into the runtime interface format and
// applies an assert.
// NOTE: for use use only in unit tests.
func JSON(data string) map[string]interface{} {
	var j map[string]interface{}
	err := json.Unmarshal([]byte(data), &j)
	gomega.Expect(err).To(gomega.BeNil())
	return j
}
