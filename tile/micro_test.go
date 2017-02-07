package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Micro", func() {

	var micro *tile.Micro

	BeforeEach(func() {
		micro = &tile.Micro{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"lod": 4
				}`)
			err := micro.Parse(params)
			Expect(err).To(BeNil())
			Expect(micro.LOD).To(Equal(4))
		})
	})

	Describe("ParseIncludes", func() {
		It("should ensure provide xField and yField are included in the includes", func() {
			includesA := micro.ParseIncludes([]string{}, "a", "b")
			includesB := micro.ParseIncludes([]string{"a"}, "a", "b")
			includesC := micro.ParseIncludes([]string{"b"}, "a", "b")
			includesD := micro.ParseIncludes([]string{"a", "b"}, "a", "b")
			Expect(includesA).To(Equal([]string{"a", "b"}))
			Expect(includesB).To(Equal([]string{"a", "b"}))
			Expect(includesC).To(Equal([]string{"b", "a"}))
			Expect(includesD).To(Equal([]string{"a", "b"}))
		})
	})

	Describe("Encode", func() {
		It("should encode results properly", func() {
			correct := []byte(`{"hits":null,"offsets":[0,8,8,8],"points":[1,1]}`)
			params := JSON(
				`{
					"lod": 1
				}`)
			err := micro.Parse(params)
			Expect(err).To(BeNil())
			bytes, err := micro.Encode(nil, []float32{1, 1})
			Expect(err).To(BeNil())
			Expect(bytes).To(Equal(correct))
		})
	})

})
