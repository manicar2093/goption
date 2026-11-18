package goption_test

import (
	"time"

	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	Describe("UnmarshalText", func() {
		It("int type", func() {
			var (
				expectedData       = 100
				expectedDataString = "100"
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[int]
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal(expectedData))
		})

		It("float type", func() {
			var (
				expectedData       = 100.00
				expectedDataString = "100.00"
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[float64]
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal(expectedData))
		})

		It("string type", func() {
			var (
				expectedData       = "hello"
				expectedDataString = "hello"
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[string]
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal(expectedData))
		})

		It("uuid type", func() {
			var (
				expectedData       = "1e2dd2c6-364b-4171-a906-554754eda276"
				expectedDataString = "1e2dd2c6-364b-4171-a906-554754eda276"
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[string]
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal(expectedData))
		})

		It("slice type", func() {
			var (
				expectedData = `["hello", "world"]`
				jsonData     = []byte(expectedData)
				holder       goption.Optional[[]string]
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal([]string{"hello", "world"}))
		})

		It("empty string type", func() {
			var (
				expectedDataString = ""
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[string]
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.IsPresent()).To(BeFalse())
		})

		DescribeTable("time type", func(expectedDataString []byte) {
			var holder goption.Optional[time.Time]

			err := holder.UnmarshalText(expectedDataString)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.IsPresent()).To(BeFalse())
		},
			Entry("null time", []byte("null")),
			Entry("empty time", []byte("")),
		)

		Context("handle bool values from relaxed strings", func() {
			DescribeTable("if string data is in accepted types and optional is bool", func(expectedValue string, expectedResult bool) {
				var holder goption.Optional[bool]
				err := holder.UnmarshalText([]byte(expectedValue))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeTrue())
				Expect(holder.MustGet()).To(Equal(expectedResult))
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

		It("number in text", func() {
			var (
				expectedDataString = "5555555555"
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[string]
			)
			err := holder.UnmarshalText(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.IsPresent()).To(BeTrue())
			Expect(holder.MustGet()).To(Equal("5555555555"))
		})
	})
})
