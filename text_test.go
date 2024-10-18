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
