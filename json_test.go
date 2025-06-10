package goption_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ExampleStuct struct {
	Name     goption.Optional[string] `json:"name"`
	LastName goption.Optional[string] `json:"last_name,omitempty"`
	Age      goption.Optional[int]    `json:"age"`
}

func asJsonString(v any) string {
	return fmt.Sprintf(
		`"%v"`,
		v,
	)
}

var _ = Describe("Json", func() {

	Describe("UnmarshalJSON", func() {

		It("generates optional from json number", func() {
			var (
				expectedNameData       = 10000
				expectedNameDataString = "10000"
				jsonData               = []byte(expectedNameDataString)
				holder                 = goption.Empty[int]()
			)
			err := holder.UnmarshalJSON(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		})

		It("generates optional from json float", func() {
			var (
				expectedNameData       = 100.00
				expectedNameDataString = "100.00"
				jsonData               = []byte(expectedNameDataString)
				holder                 = goption.Empty[float64]()
			)
			err := holder.UnmarshalJSON(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		})

		DescribeTable("generates optional from json string", func(expectedNameData any) {
			var (
				jsonData = []byte(asJsonString(expectedNameData))
				holder   = goption.Empty[string]()
			)
			err := holder.UnmarshalJSON(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.Get()).To(Equal(expectedNameData))
		},
			Entry("one line string", "a name"),
			Entry("multiline strings", `a name
a name
a name

a name`),
			Entry("multiline strings with carriage return", "\r\na name\r\n"),
		)

		When("is as null", func() {
			DescribeTable("creates an empty optional", func(jsonData []byte) {
				var (
					holder = goption.Empty[string]()
				)
				err := holder.UnmarshalJSON(jsonData)

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeFalse())
			},
				Entry("null as string", []byte(`"null"`)),
				Entry("native null", []byte(`null`)),
			)

			It("handles custom types", func() {
				type testType time.Time
				var (
					holder = goption.Empty[testType]()
				)
				err := holder.UnmarshalJSON([]byte(`"null"`))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeFalse())
			})

		})

		When("is zero", func() {
			It("creates an empty optional", func() {
				var (
					jsonData = []byte(`""`)
					holder   = goption.Empty[string]()
				)
				err := holder.UnmarshalJSON(jsonData)

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeFalse())
			})
		})

		Context("handle bool values from strings", func() {
			DescribeTable("if string data is in accepted types and optional is bool", func(expectedValue string, expectedResult bool) {
				var (
					holder = goption.Empty[bool]()
				)
				err := holder.UnmarshalJSON([]byte(expectedValue))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeTrue())
			},
				Entry("value true", `"true"`, true),
				Entry("value 1", `"1"`, true),
				Entry("value on", `"on"`, true),
				Entry("value yes", `"yes"`, true),
				Entry("value false", `"false"`, false),
				Entry("value 0", `"0"`, false),
				Entry("value off", `"off"`, false),
				Entry("value no", `"no"`, false),
			)

			DescribeTable("if string data is in accepted types and optional is pointer bool", func(expectedValue string, expectedResult bool) {
				var (
					holder = goption.Empty[*bool]()
				)
				err := holder.UnmarshalJSON([]byte(expectedValue))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(Equal(expectedResult))
			},
				Entry("value true", `"true"`, true),
				Entry("value 1", `"1"`, true),
				Entry("value on", `"on"`, true),
				Entry("value yes", `"yes"`, true),
				Entry("value false", `"false"`, false),
				Entry("value 0", `"0"`, false),
				Entry("value off", `"off"`, false),
				Entry("value no", `"no"`, false),
			)
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
