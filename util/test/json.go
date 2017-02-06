package test

import (
	"encoding/json"
	. "github.com/onsi/gomega"
)

func JSON(data string) map[string]interface{} {
	var j map[string]interface{}
	err := json.Unmarshal([]byte(data), &j)
	Expect(err).To(BeNil())
	return j
}
