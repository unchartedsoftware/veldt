package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bivariate", func() {
	bv := &tile.Bivariate{}
	bv2 := &tile.Bivariate{}

	params := make(map[string]interface{})
	params["xField"] = "a"
	params["yField"] = "b"
	params["left"] = 1.0
	params["right"] = 1.0
	params["bottom"] = 1.0
	params["top"] = 1.0

	params_fail := make(map[string]interface{})

	It("should set fields", func() {
		ok := bv.Parse(params)
		Expect(ok).To(BeNil())
		Expect(bv.XField).To(Equal(params["xField"]))
		Expect(bv.YField).To(Equal(params["yField"]))
		Expect(bv.Left).To(Equal(params["left"]))
		Expect(bv.Right).To(Equal(params["right"]))
		Expect(bv.Bottom).To(Equal(params["bottom"]))
		Expect(bv.Top).To(Equal(params["top"]))
	})

	It("should fail on wrong input", func() {
		ok := bv2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
