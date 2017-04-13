package salt

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSalt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Salt Suite")
}
