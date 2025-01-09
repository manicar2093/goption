package goption_test

import (
	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	Describe("UnmarshalText", func() {
		It("int type", func() {
			var (
				expectedNameData       = 100
				expectedNameDataString = "100"
				jsonData               = []byte(expectedNameDataString)
				holder                 = goption.Empty[int]()
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		})

		It("float type", func() {
			var (
				expectedNameData       = 100.00
				expectedNameDataString = "100.00"
				jsonData               = []byte(expectedNameDataString)
				holder                 = goption.Empty[float64]()
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		})

		It("string type", func() {
			var (
				expectedNameData       = "hello"
				expectedNameDataString = "hello"
				jsonData               = []byte(expectedNameDataString)
				holder                 = goption.Empty[string]()
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		})

		It("uuid type", func() {
			var (
				expectedNameData       = "1e2dd2c6-364b-4171-a906-554754eda276"
				expectedNameDataString = "1e2dd2c6-364b-4171-a906-554754eda276"
				jsonData               = []byte(expectedNameDataString)
				holder                 = goption.Empty[string]()
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		})

		It("slice type", func() {
			var (
				expectedNameData = `["hello", "world"]`
				jsonData         = []byte(expectedNameData)
				holder           = goption.Empty[[]string]()
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal([]string{"hello", "world"}))
		})

		It("slices objects type", func() {
			type test struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			var (
				expectedNameData = `[{"name":"hello","age":20},{"name":"hello2","age":30}]`
				jsonData         = []byte(expectedNameData)
				holder           = goption.Empty[[]test]()
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal([]test{
				{
					Name: "hello",
					Age:  20,
				}, {
					Name: "hello2",
					Age:  30,
				},
			}))
		})

		It("empty string type", func() {
			var (
				expectedNameDataString = ""
				jsonData               = []byte(expectedNameDataString)
				holder                 = goption.Empty[string]()
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.IsPresent()).To(BeFalse())
		})
	})
})
