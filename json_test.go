package goption_test

import (
	"encoding/json"
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

	Describe("MarshalJson", func() {
		It("creates json from optional", func() {
			var (
				expectedValue    = "my name"
				expectedJson     = `"my name"`
				expectedOptional = goption.Of(expectedValue)
			)

			got, err := json.Marshal(expectedOptional)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(got)).To(Equal(expectedJson))
		})

		It("creates json when optional is in a struct", func() {
			var (
				expectedName     = "expectedName"
				expectedLastName = "expectedLastName"
				expectedAge      = 40
				expectedStruct   = ExampleStuct{
					Name:     goption.Of(expectedName),
					LastName: goption.Of(expectedLastName),
					Age:      goption.Of(expectedAge),
				}
			)

			got, err := json.Marshal(expectedStruct)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(
				MatchJSON(
					fmt.Sprintf(`{"name":"%v",
			"last_name":"%v",
			"age":%v}`, expectedName, expectedLastName, expectedAge),
				),
			)
		})

		When("optional has omitempty in a struct", func() {
			It("send it as null in json", func() {
				var (
					expectedName   = "expectedName"
					expectedAge    = 40
					expectedStruct = ExampleStuct{
						Name: goption.Of(expectedName),
						Age:  goption.Of(expectedAge),
					}
				)

				got, err := json.Marshal(expectedStruct)

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(
					MatchJSON(
						fmt.Sprintf(`{"name":"%v",
						"last_name":null,
				"age":%v}`, expectedName, expectedAge),
					),
				)
			})
		})
	})

})
