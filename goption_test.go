package goption_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/goption"
)

func toPointer[T any](value T) *T {
	var pointer T = value
	return &pointer
}

var _ = Describe("Goption", func() {

	Describe("Empty", func() {
		It("creates an empty Optional", func() {
			got := goption.Empty[string]()
			Expect(got.IsPresent()).To(BeFalse())
			Expect(got.IsZero()).To(BeTrue())
		})
	})

	Describe("Of", func() {
		DescribeTable("creates an Optional with given data", func(expectedValue interface{}) {
			got := goption.Of(expectedValue)

			Expect(got.IsPresent()).To(BeTrue())
			Expect(got.IsZero()).To(BeFalse())
		},
			Entry("no pointer string value", "hello"),
			Entry("no pointer int value", int(6)),
			Entry("no pointer int8 value", int8(6)),
			Entry("no pointer uint16 value", uint16(6)),
			Entry("no pointer uint32 value", uint32(6)),
			Entry("no pointer uint64 value", uint64(6)),
			Entry("no pointer uint8 value", uint8(6)),
			Entry("no pointer uint16 value", uint16(6)),
			Entry("no pointer uint32 value", uint32(6)),
			Entry("no pointer uint64 value", uint64(6)),
			Entry("no pointer complex64 value", complex64(6)),
			Entry("no pointer complex128 value", complex128(6)),
			Entry("no pointer struct value", struct{ Name string }{Name: "name"}),
			Entry("pointer string value", toPointer("hello")),
			Entry("pointer int value", toPointer(int(6))),
			Entry("pointer int8 value", toPointer(int8(6))),
			Entry("pointer uint16 value", toPointer(uint16(6))),
			Entry("pointer uint32 value", toPointer(uint32(6))),
			Entry("pointer uint64 value", toPointer(uint64(6))),
			Entry("pointer uint8 value", toPointer(uint8(6))),
			Entry("pointer uint16 value", toPointer(uint16(6))),
			Entry("pointer uint32 value", toPointer(uint32(6))),
			Entry("pointer uint64 value", toPointer(uint64(6))),
			Entry("pointer complex64 value", toPointer(complex64(6))),
			Entry("pointer complex128 value", toPointer(complex128(6))),
			Entry("pointer struct value", toPointer(struct{ Name string }{Name: "name"})),
		)
	})

	Describe("IsPresent", func() {
		DescribeTable("identify if optional has a valid data", func(expectedValue interface{}, isPresent bool) {
			got := goption.Of(expectedValue)
			Expect(got.IsPresent()).To(Equal(isPresent))
			Expect(got.IsZero()).To(Equal(!isPresent))
		},
			Entry("data is nil", nil, false),
			Entry("filled pointer", toPointer(struct{ Name string }{Name: "name"}), true),
			Entry("empty struct", struct{ Name string }{}, false),
			Entry("empty string", "", false),
			Entry("filled string", "full string", true),
			Entry("empty int", int(0), false),
			Entry("filled int", int(1), true),
			Entry("empty optional", goption.Optional[string]{}, false),
			Entry("empty slice", []string{}, false),
			Entry("empty slice of pointers", []*string{}, false),
		)
	})

	Describe("OrElseError", func() {
		When("optional is empty", func() {
			It("returns given error", func() {
				var (
					opt           = goption.Empty[string]()
					expectedError = errors.New("expected error")
				)

				got, err := opt.OrElseError(expectedError)

				Expect(got).To(BeEmpty())
				Expect(err).To(Equal(expectedError))
			})
		})

		When("optional has a valid value", func() {
			It("returns it", func() {
				var (
					expectedValue = "expectedValue"
					opt           = goption.Of(expectedValue)
					expectedError = errors.New("expected error")
				)

				got, err := opt.OrElseError(expectedError)

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal(expectedValue))
			})
		})
	})

	Describe("OrElse", func() {
		When("optional is empty", func() {
			It("returns given data", func() {
				var (
					opt               = goption.Empty[string]()
					expectedOtherData = "other expected data"
				)

				got := opt.OrElse(expectedOtherData)

				Expect(got).To(Equal(expectedOtherData))
			})
		})

		When("optional has valid data", func() {
			It("returns it", func() {
				var (
					expectedValue     = "expectedValue"
					opt               = goption.Of(expectedValue)
					expectedOtherData = "other expected data"
				)

				got := opt.OrElse(expectedOtherData)

				Expect(got).To(Equal(expectedValue))
			})
		})
	})

	Describe("Get", func() {
		When("optional has a valida data", func() {
			It("returns it", func() {
				var (
					expectedValue string = "a valid value"
					opt                  = goption.Of(expectedValue)
				)

				got, _ := opt.Get()

				Expect(got).To(Equal(expectedValue))
			})
		})

		When("optional is empty", func() {
			It("returns a goption.ErrNoSuchElement", func() {
				got, err := goption.Empty[string]().Get()

				Expect(got).To(BeEmpty())
				Expect(err).To(MatchError(goption.ErrNoSuchElement))
			})
		})
	})

	Describe("MustGet", func() {
		When("optional has no data", func() {
			It("panics with ErrNoSuchElement", func() {
				var opt = goption.Empty[string]()

				Expect(func() {
					opt.MustGet()
				}).To(PanicWith(goption.ErrNoSuchElement))
			})
		})

		When("optional has valid data", func() {
			It("returns it", func() {

				var (
					expectedValue = "expectedValue"
					opt           = goption.Of(expectedValue)
				)

				Expect(opt.MustGet()).To(Equal(expectedValue))
			})
		})
	})

})
