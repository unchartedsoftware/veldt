package param_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestParam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Param Suite")
}
