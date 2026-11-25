package goption_test

import (
	"encoding/json"
	"time"

	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ExampleStruct struct {
	Name     goption.Optional[string] `json:"name"`
	LastName goption.Optional[string] `json:"last_name,omitempty"`
	Age      goption.Optional[int]    `json:"age"`
}

type UserPayload struct {
	Username goption.Optional[string] `json:"username,omitempty"`
	Bio      goption.Optional[string] `json:"bio,omitempty"`
	Age      goption.Optional[int]    `json:"age,omitempty"`
	Active   goption.Optional[bool]   `json:"active,omitempty"`
}

type UserPayloadZero struct {
	Username goption.Zero[string] `json:"username,omitempty"`
	Bio      goption.Zero[string] `json:"bio,omitempty"`
	Age      goption.Zero[int]    `json:"age,omitempty"`
	Active   goption.Zero[bool]   `json:"active,omitempty"`
}

type OmitZeroStruct struct {
	Title goption.Optional[string] `json:"title,omitzero"`
	Count goption.Optional[int]    `json:"count,omitzero"`
}

type OmitZeroStructZero struct {
	Title goption.Zero[string] `json:"title,omitzero"`
	Count goption.Zero[int]    `json:"count,omitzero"`
}

func asJsonString(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

var _ = Describe("Json", func() {
	It("preserves defined string JSON behavior", func() {
		var opt goption.Optional[definedTextString]
		var zero goption.Zero[definedTextString]
		Expect(json.Unmarshal([]byte(`"hello"`), &opt)).To(Succeed())
		Expect(json.Unmarshal([]byte(`"hello"`), &zero)).To(Succeed())
		Expect(opt.MustGet()).To(Equal(definedTextString("hello")))
		Expect(zero.MustGet()).To(Equal(definedTextString("hello")))
		optJSON, err := json.Marshal(opt)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(optJSON)).To(Equal(`"hello"`))
		zeroJSON, err := json.Marshal(zero)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(zeroJSON)).To(Equal(`"hello"`))
		Expect(json.Unmarshal([]byte(`""`), &opt)).To(Succeed())
		Expect(json.Unmarshal([]byte(`""`), &zero)).To(Succeed())
		Expect(opt.IsPresent()).To(BeTrue())
		Expect(zero.IsPresent()).To(BeFalse())
		Expect(json.Unmarshal([]byte(`null`), &opt)).To(Succeed())
		Expect(json.Unmarshal([]byte(`null`), &zero)).To(Succeed())
		Expect(opt.IsPresent()).To(BeFalse())
		Expect(zero.IsPresent()).To(BeFalse())
	})

	Describe("UnmarshalJSON", func() {

		It("generates optional from json number", func() {
			var (
				expectedData       = 10000
				expectedDataString = "10000"
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[int]
			)
			err := holder.UnmarshalJSON(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal(expectedData))
		})

		It("generates optional from json float", func() {
			var (
				expectedData       = 100.00
				expectedDataString = "100.00"
				jsonData           = []byte(expectedDataString)
				holder             goption.Optional[float64]
			)
			err := holder.UnmarshalJSON(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal(expectedData))
		})

		DescribeTable("generates optional from json string", func(expectedNameData any) {
			var (
				jsonData = []byte(asJsonString(expectedNameData))
				holder   goption.Optional[string]
			)
			err := holder.UnmarshalJSON(jsonData)

			Expect(err).ToNot(HaveOccurred())
			Expect(holder.MustGet()).To(Equal(expectedNameData))
		},
			Entry("one line string", "a name"),
			Entry("multiline strings", "a name\na name\na name\n\na name"),
			Entry("multiline strings with carriage return", "\r\na name\r\n"),
		)

		When("is null", func() {
			It("creates an empty optional from literal null", func() {
				var holder = goption.Some("initial")
				err := holder.UnmarshalJSON([]byte(`null`))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeFalse())
			})

			It("handles custom types", func() {
				type testType time.Time
				var holder goption.Optional[testType]
				err := holder.UnmarshalJSON([]byte(`null`))

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeFalse())
			})
		})

		When("is zero value in JSON", func() {
			It("Optional[string] keeps empty string as valid present", func() {
				var (
					jsonData = []byte(`""`)
					holder   goption.Optional[string]
				)
				err := holder.UnmarshalJSON(jsonData)

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeTrue())
				Expect(holder.MustGet()).To(Equal(""))
			})

			It("Zero[string] marks empty string as absent", func() {
				var (
					jsonData = []byte(`""`)
					holder   goption.Zero[string]
				)
				err := holder.UnmarshalJSON(jsonData)

				Expect(err).ToNot(HaveOccurred())
				Expect(holder.IsPresent()).To(BeFalse())
			})
		})

		Context("Unmarshal structs with omitted vs null vs zero fields", func() {
			It("unmarshals partial JSON with omitted fields", func() {
				// bio and age are completely omitted from JSON
				jsonData := []byte(`{"username":"john_doe","active":false}`)

				var payload UserPayload
				err := json.Unmarshal(jsonData, &payload)
				Expect(err).ToNot(HaveOccurred())

				// username was present
				Expect(payload.Username.IsPresent()).To(BeTrue())
				Expect(payload.Username.MustGet()).To(Equal("john_doe"))

				// active was present with zero value (false)
				Expect(payload.Active.IsPresent()).To(BeTrue())
				Expect(payload.Active.MustGet()).To(Equal(false))

				// bio and age were omitted -> invalid/absent
				Expect(payload.Bio.IsPresent()).To(BeFalse())
				Expect(payload.Age.IsPresent()).To(BeFalse())
			})

			It("unmarshals explicit null vs zero value for Optional[T]", func() {
				jsonData := []byte(`{"username":"","bio":null,"age":0}`)

				var payload UserPayload
				err := json.Unmarshal(jsonData, &payload)
				Expect(err).ToNot(HaveOccurred())

				// username is valid empty string
				Expect(payload.Username.IsPresent()).To(BeTrue())
				Expect(payload.Username.MustGet()).To(Equal(""))

				// bio is explicit null
				Expect(payload.Bio.IsPresent()).To(BeFalse())

				// age is valid 0
				Expect(payload.Age.IsPresent()).To(BeTrue())
				Expect(payload.Age.MustGet()).To(Equal(0))
			})

			It("unmarshals explicit null vs zero value for Zero[T]", func() {
				jsonData := []byte(`{"username":"alice","bio":"","age":0}`)

				var payload UserPayloadZero
				err := json.Unmarshal(jsonData, &payload)
				Expect(err).ToNot(HaveOccurred())

				// username is non-zero string -> present
				Expect(payload.Username.IsPresent()).To(BeTrue())
				Expect(payload.Username.MustGet()).To(Equal("alice"))

				// bio is zero string -> absent in Zero[T]
				Expect(payload.Bio.IsPresent()).To(BeFalse())

				// age is 0 -> absent in Zero[T]
				Expect(payload.Age.IsPresent()).To(BeFalse())
			})
		})
	})

	Describe("MarshalJSON", func() {
		It("creates json from valid Optional", func() {
			expectedOptional := goption.Some("my name")

			got, err := json.Marshal(expectedOptional)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(got)).To(Equal(`"my name"`))
		})

		It("creates json when optional is in a struct", func() {
			expectedStruct := ExampleStruct{
				Name:     goption.Some("expectedName"),
				LastName: goption.Some("expectedLastName"),
				Age:      goption.Some(40),
			}

			got, err := json.Marshal(expectedStruct)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(
				MatchJSON(`{"name":"expectedName","last_name":"expectedLastName","age":40}`),
			)
		})

		Context("JSON tags: omitempty behavior", func() {
			It("marshals invalid Optional fields tagged with omitempty as null", func() {
				payload := UserPayload{
					Username: goption.Some("john"),
					// Bio is uninitialized (invalid Optional)
					Age: goption.Some(0), // zero value
				}

				got, err := json.Marshal(payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(got).To(MatchJSON(`{
					"username": "john",
					"bio": null,
					"age": 0,
					"active": null
				}`))
			})

			It("marshals Zero fields with omitempty as null when absent/zero", func() {
				payload := UserPayloadZero{
					Username: goption.SomeZero("john"),
					Bio:      goption.SomeZero(""), // zero value -> absent
					Age:      goption.SomeZero(0),  // zero value -> absent
				}

				got, err := json.Marshal(payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(got).To(MatchJSON(`{
					"username": "john",
					"bio": null,
					"age": null,
					"active": null
				}`))
			})
		})

		Context("JSON tags: omitzero (Go 1.24+) behavior via IsZero()", func() {
			It("omits invalid Optional fields when tagged with omitzero", func() {
				payload := OmitZeroStruct{
					Title: goption.Some("Hello"),
					// Count is uninitialized (IsZero() == true)
				}

				got, err := json.Marshal(payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(got).To(MatchJSON(`{"title":"Hello"}`))
			})

			It("preserves zero values in Optional fields when tagged with omitzero", func() {
				payload := OmitZeroStruct{
					Title: goption.Some(""), // valid zero string (IsZero() == false)
					Count: goption.Some(0),  // valid zero int (IsZero() == false)
				}

				got, err := json.Marshal(payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(got).To(MatchJSON(`{"title":"","count":0}`))
			})

			It("omits empty/zero fields in Zero[T] when tagged with omitzero", func() {
				payload := OmitZeroStructZero{
					Title: goption.SomeZero(""), // zero string -> IsZero() == true
					Count: goption.SomeZero(25), // non-zero int -> IsZero() == false
				}

				got, err := json.Marshal(payload)
				Expect(err).ToNot(HaveOccurred())

				Expect(got).To(MatchJSON(`{"count":25}`))
			})
		})
	})
})
