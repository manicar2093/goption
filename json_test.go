package goption_test

import (
	"fmt"

	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ExampleStuct struct {
	Name     goption.Optional[string] `json:"name"`
	LastName goption.Optional[string] `json:"last_name,omitempty"`
	Age      goption.Optional[int]    `json:"age"`
}

var _ = Describe("Json", func() {

	Describe("UnmarshalJSON", func() {

		It("generates optional from json string", func() {
			var (
				expectedNameData = "a name"
				jsonData         = []byte(fmt.Sprintf(
					`"%v"`,
					expectedNameData,
				))
				holder = goption.Empty[string]()
			)
			err := holder.UnmarshalJSON(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		})

		When("is as null", func() {
			It("creates an empty optional", func() {
				var (
					jsonData = []byte(`"null"`)
					holder   = goption.Empty[string]()
				)
				err := holder.UnmarshalJSON(jsonData)

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeFalse())
			})
		})
	})

})
