package goption_test

import (
	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
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

		DescribeTable("time type", func(expectedNameDataString []byte) {
			var holder = goption.Empty[time.Time]()

			err := holder.UnmarshalText(expectedNameDataString)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.IsPresent()).To(BeFalse())
		},
			Entry("Zero time", []byte("0001-01-01T00:00:00Z")),
			Entry("null time", []byte("null")),
			Entry("empty time", []byte("")),
		)

		Context("handle bool values from strings", func() {
			DescribeTable("if string data is in accepted types and optional is bool", func(expectedValue string, expectedResult bool) {
				var (
					holder = goption.Empty[bool]()
				)
				err := holder.UnmarshalText([]byte(expectedValue))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeTrue())
			},
				Entry("value true", "true", true),
				Entry("value 1", "1", true),
				Entry("value on", "on", true),
				Entry("value yes", "yes", true),
				Entry("value false", "false", false),
				Entry("value 0", "0", false),
				Entry("value off", "off", false),
				Entry("value no", "no", false),
			)

			DescribeTable("if string data is in accepted types and optional is pointer bool", func(expectedValue string, expectedResult bool) {
				var (
					holder = goption.Empty[*bool]()
				)
				err := holder.UnmarshalText([]byte(expectedValue))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(Equal(expectedResult))
			},
				Entry("value true", "true", true),
				Entry("value 1", "1", true),
				Entry("value on", "on", true),
				Entry("value yes", "yes", true),
				Entry("value false", "false", false),
				Entry("value 0", "0", false),
				Entry("value off", "off", false),
				Entry("value no", "no", false),
			)

		})
	})
})
