package goption_test

import (
	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	Describe("UnmarshalText", func() {
		It("float type", func() {
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
	})
})
