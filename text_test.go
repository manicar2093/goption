package goption_test

import (
	"strings"
	"time"

	"github.com/manicar2093/goption"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type definedTextString string
type definedTextBool bool

type customTextString string

func (s customTextString) MarshalText() ([]byte, error) {
	return []byte("custom:" + string(s)), nil
}

func (s *customTextString) UnmarshalText(text []byte) error {
	*s = customTextString(strings.ToUpper(string(text)))
	return nil
}

var _ = Describe("Text", func() {
	It("handles relaxed defined boolean text and preserves presence", func() {
		var opt goption.Optional[definedTextBool]
		var zero goption.Zero[definedTextBool]
		for _, tc := range []struct {
			text  string
			value bool
		}{
			{"true", true}, {"1", true}, {"on", true}, {"yes", true},
			{"false", false}, {"0", false}, {"off", false}, {"no", false},
			{" YES ", true}, {"  ", false},
		} {
			Expect(opt.UnmarshalText([]byte(tc.text))).To(Succeed())
			Expect(zero.UnmarshalText([]byte(tc.text))).To(Succeed())
			Expect(opt.IsPresent()).To(BeTrue())
			Expect(opt.MustGet()).To(Equal(definedTextBool(tc.value)))
			Expect(zero.IsPresent()).To(Equal(tc.value))
			optText, err := opt.MarshalText()
			Expect(err).NotTo(HaveOccurred())
			var roundTrip goption.Optional[definedTextBool]
			Expect(roundTrip.UnmarshalText(optText)).To(Succeed())
			Expect(roundTrip.MustGet()).To(Equal(definedTextBool(tc.value)))
			zeroText, err := zero.MarshalText()
			Expect(err).NotTo(HaveOccurred())
			var zeroRoundTrip goption.Zero[definedTextBool]
			Expect(zeroRoundTrip.UnmarshalText(zeroText)).To(Succeed())
			Expect(zeroRoundTrip.IsPresent()).To(Equal(tc.value))
		}
		Expect(opt.UnmarshalText([]byte("invalid"))).NotTo(Succeed())
		Expect(zero.UnmarshalText([]byte("invalid"))).NotTo(Succeed())
		for _, input := range []string{"", "null"} {
			Expect(opt.UnmarshalText([]byte("true"))).To(Succeed())
			Expect(zero.UnmarshalText([]byte("true"))).To(Succeed())
			Expect(opt.UnmarshalText([]byte(input))).To(Succeed())
			Expect(zero.UnmarshalText([]byte(input))).To(Succeed())
			Expect(opt.IsPresent()).To(BeFalse())
			Expect(zero.IsPresent()).To(BeFalse())
		}
	})

	It("round trips defined strings and clears absent inputs", func() {
		var opt goption.Optional[definedTextString]
		var zero goption.Zero[definedTextString]
		for _, input := range []string{"hello", "555"} {
			Expect(opt.UnmarshalText([]byte(input))).To(Succeed())
			Expect(zero.UnmarshalText([]byte(input))).To(Succeed())
			Expect(opt.MustGet()).To(Equal(definedTextString(input)))
			Expect(zero.MustGet()).To(Equal(definedTextString(input)))
			optText, err := opt.MarshalText()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(optText)).To(Equal(input))
			zeroText, err := zero.MarshalText()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(zeroText)).To(Equal(input))
		}
		for _, input := range []string{"", "null"} {
			Expect(opt.UnmarshalText([]byte("hello"))).To(Succeed())
			Expect(zero.UnmarshalText([]byte("hello"))).To(Succeed())
			Expect(opt.UnmarshalText([]byte(input))).To(Succeed())
			Expect(zero.UnmarshalText([]byte(input))).To(Succeed())
			Expect(opt.IsPresent()).To(BeFalse())
			Expect(zero.IsPresent()).To(BeFalse())
		}
	})

	It("prefers custom text interfaces for defined strings", func() {
		var opt goption.Optional[customTextString]
		var zero goption.Zero[customTextString]
		Expect(opt.UnmarshalText([]byte("hello"))).To(Succeed())
		Expect(zero.UnmarshalText([]byte("hello"))).To(Succeed())
		Expect(opt.MustGet()).To(Equal(customTextString("HELLO")))
		Expect(zero.MustGet()).To(Equal(customTextString("HELLO")))
		optText, err := opt.MarshalText()
		Expect(err).NotTo(HaveOccurred())
		Expect(string(optText)).To(Equal("custom:HELLO"))
		zeroText, err := zero.MarshalText()
		Expect(err).NotTo(HaveOccurred())
		Expect(string(zeroText)).To(Equal("custom:HELLO"))
	})

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
