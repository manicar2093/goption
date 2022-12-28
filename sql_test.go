package goption_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/goption"
)

var _ = Describe("Sql", func() {

	Describe("Scan", func() {
		It("assigns given data to optional", func() {
			opt := goption.Empty[string]()

			Expect(opt.Scan("Hello!")).To(Succeed())
			Expect(opt.IsPresent()).To(BeTrue())
		})

		When("data is not valid", func() {
			It("create a empty optional from empty src", func() {
				opt := goption.Empty[string]()

				Expect(opt.Scan("")).To(Succeed())
				Expect(opt.IsPresent()).To(BeFalse())
			})

			It("create a empty optional from nil", func() {
				opt := goption.Empty[*string]()

				Expect(opt.Scan(nil)).To(Succeed())
				Expect(opt.IsPresent()).To(BeFalse())
			})
		})
	})

})
