package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bivariate", func() {
	eq := &tile.Bivariate{}
	eq2 := &tile.Bivariate{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	params["xField"] = "a"
	params["yField"] = "b"
	params["left"] = 1.0
	params["right"] = 1.0
	params["bottom"] = 1.0
	params["top"] = 1.0

	params_fail := make(map[string]interface{})

	It("should set fields", func() {
		ok := eq.Parse(params)
		Expect(ok).To(BeNil())
		Expect(eq.XField).To(Equal(params["xField"]))
		Expect(eq.YField).To(Equal(params["yField"]))
		Expect(eq.Left).To(Equal(params["left"]))
		Expect(eq.Right).To(Equal(params["right"]))
		Expect(eq.Bottom).To(Equal(params["bottom"]))
		Expect(eq.Top).To(Equal(params["top"]))
	})

	It("should fail on wrong input", func() {
		ok := eq2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
